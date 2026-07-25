package setup

import (
	"net/http"
	"strings"
)

// FirstRunGate middleware redirects requests to /setup if no user exists in the system.
// Static assets (/static), health endpoint (/health), and /setup routes are excluded.
func FirstRunGate(svc *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path

			// Allow static assets, health, and setup endpoints through
			if strings.HasPrefix(path, "/static") || path == "/health" || strings.HasPrefix(path, "/setup") {
				next.ServeHTTP(w, r)
				return
			}

			isFirstRun, err := svc.IsFirstRun(r.Context())
			if err == nil && isFirstRun {
				http.Redirect(w, r, "/setup", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
