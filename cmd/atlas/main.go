// Atlas — entry point.
//
// What: Wires all dependencies and starts the HTTP server.
// Why:  main.go is the only place where concrete types are constructed and
//       injected into their dependents. Nothing else in the codebase
//       instantiates its own dependencies (no global vars, no init()).
// How:
//  1. Load config (env vars via Viper)
//  2. Build logger (Zap — console in dev, JSON in prod)
//  3. Create /data directories if they don't exist
//  4. Register routes on a Chi router
//  5. Start HTTP server — block until shutdown signal
//
// Phase 0: Server starts with /health only.
// Phase 1: Database connection added here.
// Phase 2: Auth middleware and /setup route added here.
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

	"github.com/Suke2004/atlas-go/internal/config"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/health"
	"github.com/Suke2004/atlas-go/internal/logger"
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
	// Ensure /data subdirectories exist at startup so later phases
	// (DB, uploads) never fail because the path is missing.
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

	// ── 5. Router ──────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Core middleware — order matters.
	r.Use(middleware.RequestID)  // Attach X-Request-Id header to every request
	r.Use(middleware.RealIP)     // Trust X-Forwarded-For behind a reverse proxy
	r.Use(middleware.Recoverer)  // Recover from panics; log and return 500
	r.Use(zapRequestLogger(log)) // Structured request logging

	// Public routes (no auth required)
	r.Get("/health", health.Handler())

	// ── 5. Server ──────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in a goroutine so we can listen for shutdown signals.
	go func() {
		log.Info("Atlas listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// Block until SIGINT or SIGTERM received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Atlas shutting down gracefully...")
	// Future phases: close DB, flush caches, drain in-flight requests here.
	log.Info("Atlas stopped")
}

// zapRequestLogger returns a Chi middleware that logs every HTTP request
// using Zap structured fields.
//
// Why not use chi/middleware.Logger? It writes to a plain io.Writer.
// We want structured JSON fields (method, path, status, latency) so logs
// are machine-parseable in production.
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
