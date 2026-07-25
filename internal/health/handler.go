// Package health implements the GET /health endpoint.
//
// What: Returns a JSON object with server status, version, and uptime.
// Why:  Health checks are required by Docker, Caddy, and any load balancer.
//       They also serve as an instant sanity-check that the server started.
// How:  A single handler writes a fixed JSON response. Uptime is calculated
//       from a package-level startTime captured at import.
package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// startTime records when the process started, used to calculate uptime.
var startTime = time.Now()

// version is set at build time via -ldflags "-X health.version=v0.1.0".
// Falls back to "dev" when built without ldflags (local development).
var version = "dev"

// response is the JSON shape returned by GET /health.
type response struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// Handler returns an http.HandlerFunc that serves GET /health.
// It always returns HTTP 200 with a JSON body as long as the process is alive.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(response{
			Status:  "ok",
			Version: version,
			Uptime:  time.Since(startTime).Round(time.Second).String(),
		})
	}
}
