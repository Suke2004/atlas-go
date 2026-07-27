package analytics

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	analyticstemplates "github.com/Suke2004/atlas-go/web/templates/analytics"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
)

// Handler serves HTTP endpoints for analytics.
type Handler struct {
	service Service
	logger  *zap.Logger
}

// NewHandler constructs an analytics handler.
func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GET /analytics — render analytics dashboard page.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data, err := h.service.GetAnalyticsData(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to load analytics data", zap.Error(err))
		http.Error(w, "Failed to load analytics", http.StatusInternalServerError)
		return
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return analyticstemplates.Analytics(data).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Analytics & Insights", "/analytics", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /analytics/weekly-review — synthesize AI weekly review via HTMX/JSON.
func (h *Handler) GenerateWeeklyReview(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	review, err := h.service.GenerateWeeklyReview(r.Context(), user.ID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"review": review})
}
