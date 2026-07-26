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
	"github.com/Suke2004/atlas-go/internal/finance"
	"github.com/Suke2004/atlas-go/internal/learning"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestIntegration_FinanceLearningFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "finance_learning_integration_test.db")

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

	finRepo := finance.NewRepository(database)
	finSvc := finance.NewService(finRepo, nil)
	finHandler := finance.NewHandler(finSvc, logger)

	learnRepo := learning.NewRepository(database)
	learnSvc := learning.NewService(learnRepo)
	learnHandler := learning.NewHandler(learnSvc, logger)

	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "alex",
		DisplayName: "Alex Mercer",
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
		protected.Get("/finance", finHandler.Index)
		protected.Get("/learning", learnHandler.Index)
	})

	// 1. GET /finance
	reqFin := httptest.NewRequest("GET", "/finance", nil)
	reqFin.AddCookie(sessionCookie)
	wFin := httptest.NewRecorder()
	r.ServeHTTP(wFin, reqFin)

	if wFin.Code != http.StatusOK {
		t.Fatalf("expected GET /finance to return 200, got: %d", wFin.Code)
	}

	// 2. GET /learning
	reqLearn := httptest.NewRequest("GET", "/learning", nil)
	reqLearn.AddCookie(sessionCookie)
	wLearn := httptest.NewRecorder()
	r.ServeHTTP(wLearn, reqLearn)

	if wLearn.Code != http.StatusOK {
		t.Fatalf("expected GET /learning to return 200, got: %d", wLearn.Code)
	}
}
