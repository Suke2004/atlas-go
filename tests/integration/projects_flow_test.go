package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestIntegration_ProjectsFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "projects_integration_test.db")

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
	repo := projects.NewRepository(database)
	projectsSvc := projects.NewService(repo, nil)
	projectsHandler := projects.NewHandler(projectsSvc, logger)

	// Create user via setup service
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "alex",
		DisplayName: "Alex Mercer",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Create authenticated session token in DB
	sess, err := authSvc.CreateSession(ctx, user.ID, 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Save token in gorilla session cookie
	wInit := httptest.NewRecorder()
	reqInit := httptest.NewRequest("GET", "/", nil)
	cookieSess, _ := sessionStore.Get(reqInit)
	cookieSess.Values["token"] = sess.ID
	_ = cookieSess.Save(reqInit, wInit)
	sessionCookie := wInit.Result().Cookies()[0]

	r := chi.NewRouter()
	r.Group(func(protected chi.Router) {
		protected.Use(auth.AuthRequired(authSvc))
		protected.Get("/projects", projectsHandler.List)
		protected.Post("/projects", projectsHandler.Create)
		protected.Get("/projects/{id}", projectsHandler.Detail)
	})

	// 1. GET /projects when empty
	req := httptest.NewRequest("GET", "/projects", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected GET /projects to return 200, got: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "No matching projects") {
		t.Errorf("expected empty state message in response")
	}

	// 2. POST /projects to create a new project
	form := url.Values{}
	form.Set("name", "Atlas Personal OS")
	form.Set("description", "A self-hosted personal operating system")
	form.Set("status", "active")
	form.Set("color", "#6366f1")
	form.Set("github_url", "https://github.com/Suke2004/atlas-go")
	form.Set("tech_stack", "Go, SQLite, HTMX, TailwindCSS")

	reqPost := httptest.NewRequest("POST", "/projects", strings.NewReader(form.Encode()))
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPost.AddCookie(sessionCookie)
	wPost := httptest.NewRecorder()
	r.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusSeeOther {
		t.Fatalf("expected POST /projects to redirect with 303, got: %d", wPost.Code)
	}

	// 3. Verify project created in DB
	list, err := projectsSvc.ListProjects(ctx, user.ID, "all", "", "")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 project in DB, got: %v", list)
	}
	if list[0].Name != "Atlas Personal OS" {
		t.Errorf("expected project name 'Atlas Personal OS', got: %s", list[0].Name)
	}
	if list[0].TechStack != "Go, SQLite, HTMX, TailwindCSS" {
		t.Errorf("expected tech stack 'Go, SQLite, HTMX, TailwindCSS', got: %s", list[0].TechStack)
	}
}
