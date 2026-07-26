package learning

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	learningtemplates "github.com/Suke2004/atlas-go/web/templates/learning"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// GET /learning
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tracks, err := h.service.ListTracks(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to list learning tracks", zap.Error(err))
	}

	summary, _ := h.service.GetLearningSummary(r.Context(), user.ID)

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return learningtemplates.Learning(tracks, summary).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Skill Roadmap & Mastery", "/learning", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /learning/tracks
func (h *Handler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	input := TrackInput{
		Title:       r.FormValue("title"),
		Category:    r.FormValue("category"),
		Description: r.FormValue("description"),
	}

	_, err := h.service.CreateTrack(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to create learning track", zap.Error(err))
	}

	http.Redirect(w, r, "/learning", http.StatusSeeOther)
}

// POST /learning/sessions
func (h *Handler) AddSession(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	trackID, _ := strconv.ParseInt(r.FormValue("track_id"), 10, 64)
	minutes, _ := strconv.ParseInt(r.FormValue("duration_minutes"), 10, 64)

	_, err := h.service.AddSession(r.Context(), user.ID, trackID, minutes, r.FormValue("summary"))
	if err != nil {
		h.logger.Error("failed to add learning session", zap.Error(err))
	}

	http.Redirect(w, r, "/learning", http.StatusSeeOther)
}

// POST /learning/tracks/{id}/delete
func (h *Handler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	trackID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	_ = h.service.DeleteTrack(r.Context(), user.ID, trackID)
	http.Redirect(w, r, "/learning", http.StatusSeeOther)
}
