package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/setup"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func TestIntegration_FirstRunSetupFlow(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "setup_flow.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	logger := zap.NewNop()
	authStore := auth.NewStore("test-secret-key-32-bytes-long!", 3600)
	authSvc := auth.NewService(database, authStore)
	setupSvc := setup.NewService(database)
	setupHandler := setup.NewHandler(setupSvc, authSvc, logger)

	r := chi.NewRouter()
	r.Use(setup.FirstRunGate(setupSvc))
	r.Get("/setup", setupHandler.ShowWizard)
	r.Post("/setup", setupHandler.ProcessSetup)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Dashboard"))
	})

	// 1. Initial request to / should be redirected to /setup (FirstRunGate)
	req1 := httptest.NewRequest("GET", "/", nil)
	rec1 := httptest.NewRecorder()
	r.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect 303 for first-run gate, got %d", rec1.Code)
	}

	if loc := rec1.Header().Get("Location"); loc != "/setup" {
		t.Fatalf("expected redirect to /setup, got %s", loc)
	}

	// 2. GET /setup should render 200 OK HTML
	req2 := httptest.NewRequest("GET", "/setup", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /setup, got %d", rec2.Code)
	}

	if !strings.Contains(rec2.Body.String(), "Welcome to Atlas") {
		t.Errorf("expected setup wizard HTML body, got: %s", rec2.Body.String())
	}

	// 3. POST /setup should create first user and redirect to /setup/demo-choice
	form := url.Values{}
	form.Set("username", "owner")
	form.Set("display_name", "Owner Name")
	form.Set("password", "masterpass123")
	form.Set("timezone", "UTC")
	form.Set("theme", "dark")

	req3 := httptest.NewRequest("POST", "/setup", strings.NewReader(form.Encode()))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec3 := httptest.NewRecorder()
	r.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect after POST /setup, got %d", rec3.Code)
	}

	// 4. Verify user exists in SQLite
	count, err := database.CountUsers(context.Background())
	if err != nil || count != 1 {
		t.Fatalf("expected 1 user in DB after setup, got count %d, err: %v", count, err)
	}
}
