package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestIntegration_NotesFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "notes_integration_test.db")

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
	projRepo := projects.NewRepository(database)

	notesRepo := notes.NewRepository(database)
	notesSvc := notes.NewService(notesRepo, projRepo)
	notesHandler := notes.NewHandler(notesSvc, projRepo, logger)

	// Create user
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "sam",
		DisplayName: "Sam Vance",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Create session token
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
		protected.Get("/notes", notesHandler.List)
		protected.Get("/notes/new", notesHandler.NewForm)
		protected.Post("/notes", notesHandler.Create)
		protected.Get("/notes/{id}", notesHandler.Detail)
	})

	// 1. GET /notes when empty
	req := httptest.NewRequest("GET", "/notes", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected GET /notes to return 200, got: %d", w.Code)
	}

	// 2. POST /notes to create note
	form := url.Values{}
	form.Set("title", "ADR: Database Schema Design")
	form.Set("content", "# ADR 001\nWe use SQLite with WAL mode.")
	form.Set("tags", "adr, database")

	reqPost := httptest.NewRequest("POST", "/notes", strings.NewReader(form.Encode()))
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPost.AddCookie(sessionCookie)
	wPost := httptest.NewRecorder()
	r.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusSeeOther {
		t.Fatalf("expected POST /notes to redirect with 303, got: %d", wPost.Code)
	}

	// 3. Verify note created
	list, err := notesSvc.ListNotes(ctx, user.ID, "all", "")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 note in DB, got: %v", list)
	}
	if list[0].Note.Title != "ADR: Database Schema Design" {
		t.Errorf("expected note title 'ADR: Database Schema Design', got: %s", list[0].Note.Title)
	}
}
