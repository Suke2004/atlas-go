package settings

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	settingstemplates "github.com/Suke2004/atlas-go/web/templates/settings"
)

// Handler serves the settings page.
type Handler struct {
	service Service
	logger  *zap.Logger
}

// NewHandler constructs a settings handler.
func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GET /settings — render the settings page.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	settings := h.service.GetAll(r.Context(), user.ID)
	successMsg := r.URL.Query().Get("saved")

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return settingstemplates.Settings(user.DisplayName, user.Username, settings, successMsg).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Settings", "/settings", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /settings — persist settings from form submission.
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	// All persisted setting keys
	keys := []string{
		"theme",
		"ai_provider",
		"ai_ollama_url",
		"ai_ollama_model",
		"ai_gemini_key",
		"ai_gemini_model",
		"ai_openai_key",
		"ai_openai_model",
	}
	for _, key := range keys {
		val := r.FormValue(key)
		if val == "" {
			continue
		}
		if err := h.service.Set(r.Context(), user.ID, key, val); err != nil {
			h.logger.Error("failed to save setting", zap.String("key", key), zap.Error(err))
		}
	}

	http.Redirect(w, r, "/settings?saved=1", http.StatusSeeOther)
}
