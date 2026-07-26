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
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/setup"
	"github.com/Suke2004/atlas-go/internal/tasks"
)

func TestIntegration_TasksFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "tasks_integration_test.db")

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
	projSvc := projects.NewService(projRepo, nil)

	tasksRepo := tasks.NewRepository(database)
	tasksSvc := tasks.NewService(tasksRepo, projRepo, projSvc)
	tasksHandler := tasks.NewHandler(tasksSvc, projRepo, logger)

	// Create user
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "alex",
		DisplayName: "Alex Mercer",
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
		protected.Get("/tasks", tasksHandler.List)
		protected.Post("/tasks", tasksHandler.Create)
		protected.Post("/tasks/{id}/status", tasksHandler.UpdateStatus)
		protected.Post("/tasks/{id}/delete", tasksHandler.Delete)
	})

	// 1. GET /tasks when empty
	req := httptest.NewRequest("GET", "/tasks", nil)
	req.AddCookie(sessionCookie)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected GET /tasks to return 200, got: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "No tasks found") {
		t.Errorf("expected empty state message in response")
	}

	// 2. POST /tasks to create a task
	form := url.Values{}
	form.Set("title", "Implement Kanban Board")
	form.Set("description", "3-column Kanban view for Phase 5")
	form.Set("status", "todo")
	form.Set("priority", "high")
	form.Set("energy_level", "high")

	reqPost := httptest.NewRequest("POST", "/tasks", strings.NewReader(form.Encode()))
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqPost.AddCookie(sessionCookie)
	wPost := httptest.NewRecorder()
	r.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusSeeOther {
		t.Fatalf("expected POST /tasks to redirect with 303, got: %d", wPost.Code)
	}

	// 3. Verify task created
	list, err := tasksSvc.ListTasks(ctx, user.ID, "all", "all", "all", 0, "")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 task in DB, got: %v", list)
	}
	if list[0].Task.Title != "Implement Kanban Board" {
		t.Errorf("expected task title 'Implement Kanban Board', got: %s", list[0].Task.Title)
	}
}
