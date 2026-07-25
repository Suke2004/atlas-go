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
	"github.com/Suke2004/atlas-go/internal/logger"
	"github.com/Suke2004/atlas-go/internal/setup"
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

	// ── 6. Router ──────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(zapRequestLogger(log))

	// First-Run Gate Middleware — redirects to /setup if zero users exist
	r.Use(setup.FirstRunGate(setupSvc))

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

	// Protected routes (Phase 3+ layout, dashboard, modules)
	r.Group(func(protected chi.Router) {
		protected.Use(auth.AuthRequired(authSvc))
		protected.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := auth.GetUserFromContext(r.Context())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>Atlas</title><script src='https://cdn.tailwindcss.com'></script></head><body class='bg-slate-950 text-white p-8'><div class='max-w-2xl mx-auto bg-slate-900 border border-slate-800 p-6 rounded-xl'><h1 class='text-2xl font-bold mb-2'>Atlas Workspace</h1><p class='text-slate-400 mb-4'>Welcome back, <strong>%s</strong>!</p><form action='/logout' method='POST'><button class='px-4 py-2 bg-rose-600 rounded-lg font-medium'>Log Out</button></form></div></body></html>", user.DisplayName)
		})
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
