// Atlas — entry point.
//
// What: Wires all dependencies and starts the HTTP server.
// Why:  main.go is the only place where concrete types are constructed and
//       injected into their dependents. Nothing else in the codebase
//       instantiates its own dependencies (no global vars, no init()).
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/config"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/health"
	"github.com/Suke2004/atlas-go/internal/journal"
	"github.com/Suke2004/atlas-go/internal/logger"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/search"
	"github.com/Suke2004/atlas-go/internal/setup"
	"github.com/Suke2004/atlas-go/internal/tasks"
	dashtemplates "github.com/Suke2004/atlas-go/web/templates/dashboard"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
)

func main() {
	// ── 1. Config ──────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── 2. Logger ──────────────────────────────────────────────────────────
	log := logger.New(cfg.Env, cfg.LogLevel)
	defer log.Sync() //nolint:errcheck

	log.Info("Atlas starting",
		zap.String("env", cfg.Env),
		zap.Int("port", cfg.Port),
		zap.String("version", "dev"),
	)

	// ── 3. Data directories ────────────────────────────────────────────────
	dataDirs := []string{
		"data/db",
		"data/uploads",
		"data/backups",
		"data/logs",
	}
	for _, dir := range dataDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("failed to create data directory",
				zap.String("dir", dir),
				zap.Error(err),
			)
		}
	}

	// ── 4. Database & Migrations ───────────────────────────────────────────
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatal("failed to open database", zap.Error(err))
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		log.Fatal("failed to apply migrations", zap.Error(err))
	}
	log.Info("database migrations applied successfully", zap.String("db_path", cfg.DBPath))

	// ── 5. Services & Handlers ─────────────────────────────────────────────
	sessionStore := auth.NewStore(cfg.SessionSecret, cfg.SessionMaxAge)
	authSvc := auth.NewService(database, sessionStore)
	authHandler := auth.NewHandler(authSvc, log)

	setupSvc := setup.NewService(database)
	setupHandler := setup.NewHandler(setupSvc, authSvc, log)

	projectsRepo := projects.NewRepository(database)
	projectsSvc := projects.NewService(projectsRepo, nil)
	projectsHandler := projects.NewHandler(projectsSvc, log)

	tasksRepo := tasks.NewRepository(database)
	tasksSvc := tasks.NewService(tasksRepo, projectsRepo, projectsSvc)
	tasksHandler := tasks.NewHandler(tasksSvc, projectsRepo, log)

	notesRepo := notes.NewRepository(database)
	notesSvc := notes.NewService(notesRepo, projectsRepo)
	notesHandler := notes.NewHandler(notesSvc, projectsRepo, log)

	journalRepo := journal.NewRepository(database)
	journalSvc := journal.NewService(journalRepo, tasksSvc, notesSvc)
	journalHandler := journal.NewHandler(journalSvc, log)

	searchSvc := search.NewService(database, projectsSvc, tasksSvc, notesSvc)
	searchHandler := search.NewHandler(searchSvc, log)

	// ── 6. Router ──────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// All middleware MUST be defined before any routes on Chi Mux
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(zapRequestLogger(log))
	r.Use(setup.FirstRunGate(setupSvc))

	// Serve local static assets (CSS, JS, icons)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Public endpoints
	r.Get("/health", health.Handler())

	// First-Run Setup routes
	r.Get("/setup", setupHandler.ShowWizard)
	r.Post("/setup", setupHandler.ProcessSetup)
	r.Get("/setup/demo-choice", setupHandler.ShowDemoChoice)
	r.Post("/setup/seed", setupHandler.ProcessDemoSeed)

	// Auth routes
	r.Get("/login", authHandler.ShowLogin)
	r.Post("/login", authHandler.ProcessLogin)
	r.Post("/logout", authHandler.Logout)

	// Protected application routes
	r.Group(func(protected chi.Router) {
		protected.Use(auth.AuthRequired(authSvc))

		// Dashboard route (Root Layout Shell & Executive Command Center)
		protected.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserFromContext(r.Context())
			username := "Owner"
			var userID int64
			if user != nil {
				username = user.DisplayName
				userID = user.ID
			}

			projSummary, _ := projectsSvc.GetProjectsSummary(r.Context(), userID)
			activeProjects, _ := projectsRepo.ListProjects(r.Context(), userID)
			taskSummary, _ := tasksSvc.GetTasksSummary(r.Context(), userID)
			focusTasks, _ := tasksSvc.ListTasks(r.Context(), userID, "todo", "high", "all", 0, "")

			dashContent := dashtemplates.Dashboard(username, projSummary, activeProjects, taskSummary, focusTasks)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = layouts.Base("Dashboard", "/", username, dashContent).Render(r.Context(), w)
		})

		// Projects Module Routes
		protected.Get("/projects", projectsHandler.List)
		protected.Post("/projects", projectsHandler.Create)
		protected.Get("/projects/{id}", projectsHandler.Detail)
		protected.Post("/projects/{id}/edit", projectsHandler.Update)
		protected.Post("/projects/{id}/delete", projectsHandler.Delete)
		protected.Post("/projects/{id}/sync-github", projectsHandler.SyncGitHub)
		protected.Post("/projects/{id}/import-issues", projectsHandler.ImportIssues)
		protected.Post("/projects/{id}/milestones", projectsHandler.CreateMilestone)
		protected.Post("/projects/{id}/milestones/{milestoneID}/toggle", projectsHandler.ToggleMilestone)
		protected.Post("/projects/{id}/milestones/{milestoneID}/delete", projectsHandler.DeleteMilestone)

		// Tasks Module Routes
		protected.Get("/tasks", tasksHandler.List)
		protected.Post("/tasks", tasksHandler.Create)
		protected.Post("/tasks/{id}/status", tasksHandler.UpdateStatus)
		protected.Post("/tasks/{id}/delete", tasksHandler.Delete)

		// Knowledge Base (Notes) Module Routes
		protected.Get("/notes", notesHandler.List)
		protected.Get("/notes/new", notesHandler.NewForm)
		protected.Post("/notes", notesHandler.Create)
		protected.Get("/notes/{id}", notesHandler.Detail)
		protected.Post("/notes/{id}/edit", notesHandler.Update)
		protected.Post("/notes/{id}/autosave", notesHandler.Autosave)
		protected.Post("/notes/{id}/pin", notesHandler.TogglePin)
		protected.Post("/notes/{id}/delete", notesHandler.Delete)

		// Journal Module Routes
		protected.Get("/journal", journalHandler.Index)
		protected.Post("/journal/save", journalHandler.Save)
		protected.Post("/journal/items", journalHandler.AddItem)
		protected.Post("/journal/items/{id}/delete", journalHandler.DeleteItem)

		// Global Search API Route
		protected.Get("/api/search", searchHandler.Search)
	})

	// ── 7. Server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Atlas listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Atlas shutting down gracefully...")
	log.Info("Atlas stopped")
}

func zapRequestLogger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			log.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Duration("latency", time.Since(start)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}
