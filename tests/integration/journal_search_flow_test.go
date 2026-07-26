package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/journal"
	"github.com/Suke2004/atlas-go/internal/search"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestIntegration_JournalSearchFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "journal_search_integration_test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	ctx := context.Background()
	logger := zap.NewNop()

	sessionStore := auth.NewStore("test-secret-32-bytes-long-string!", 86400)
	authSvc := auth.NewService(database, sessionStore)
	setupSvc := setup.NewService(database)

	journalRepo := journal.NewRepository(database)
	journalSvc := journal.NewService(journalRepo, nil, nil)
	journalHandler := journal.NewHandler(journalSvc, logger)

	searchSvc := search.NewService(database, nil, nil, nil)
	searchHandler := search.NewHandler(searchSvc, logger)

	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "taylor",
		DisplayName: "Taylor Swift",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	sess, err := authSvc.CreateSession(ctx, user.ID, 24*100)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	wInit := httptest.NewRecorder()
	reqInit := httptest.NewRequest("GET", "/", nil)
	cookieSess, _ := sessionStore.Get(reqInit)
	cookieSess.Values["token"] = sess.ID
	_ = cookieSess.Save(reqInit, wInit)
	sessionCookie := wInit.Result().Cookies()[0]

	r := chi.NewRouter()
	r.Group(func(protected chi.Router) {
		protected.Use(auth.AuthRequired(authSvc))
		protected.Get("/journal", journalHandler.Index)
		protected.Get("/api/search", searchHandler.Search)
	})

	// 1. GET /journal
	req := httptest.NewRequest("GET", "/journal", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected GET /journal to return 200, got: %d", w.Code)
	}

	// 2. GET /api/search
	reqSearch := httptest.NewRequest("GET", "/api/search?q=test", nil)
	reqSearch.AddCookie(sessionCookie)
	wSearch := httptest.NewRecorder()
	r.ServeHTTP(wSearch, reqSearch)

	if wSearch.Code != http.StatusOK {
		t.Fatalf("expected GET /api/search to return 200, got: %d", wSearch.Code)
	}
}
