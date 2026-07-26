package notes

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	notetemplates "github.com/Suke2004/atlas-go/web/templates/notes"
)

type Handler struct {
	service      Service
	projectsRepo projects.Repository
	logger       *zap.Logger
}

func NewHandler(service Service, projRepo projects.Repository, logger *zap.Logger) *Handler {
	return &Handler{
		service:      service,
		projectsRepo: projRepo,
		logger:       logger,
	}
}

// GET /notes
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tagFilter := r.URL.Query().Get("tag")
	searchQuery := r.URL.Query().Get("search")
	viewMode := r.URL.Query().Get("view")
	if viewMode == "" {
		viewMode = "grid"
	}

	notesList, err := h.service.ListNotes(r.Context(), user.ID, tagFilter, searchQuery)
	if err != nil {
		h.logger.Error("failed to list notes", zap.Error(err))
		http.Error(w, "Failed to load notes", http.StatusInternalServerError)
		return
	}

	summary, err := h.service.GetNotesSummary(r.Context(), user.ID)
	if err != nil {
		h.logger.Warn("failed to calculate notes summary", zap.Error(err))
	}

	var availableProjects []db.Project
	if h.projectsRepo != nil {
		availableProjects, _ = h.projectsRepo.ListProjects(r.Context(), user.ID)
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return notetemplates.List(notesList, summary, availableProjects, tagFilter, viewMode, searchQuery, "").Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Knowledge Base", "/notes", user.DisplayName, pageContent).Render(r.Context(), w)
}

// GET /notes/new
func (h *Handler) NewForm(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templateType := r.URL.Query().Get("template")
	title, content := h.service.GetTemplateContent(templateType)

	var availableProjects []db.Project
	if h.projectsRepo != nil {
		availableProjects, _ = h.projectsRepo.ListProjects(r.Context(), user.ID)
	}

	emptyNote := notetemplates.NoteWithDetails{
		Note: db.Note{
			Title:   title,
			Content: content,
		},
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return notetemplates.Editor(emptyNote, availableProjects, true).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("New Note", "/notes", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /notes
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	var projectID int64
	if projStr := r.FormValue("project_id"); projStr != "" {
		projectID, _ = strconv.ParseInt(projStr, 10, 64)
	}

	tagsRaw := r.FormValue("tags")
	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			if strings.TrimSpace(t) != "" {
				tags = append(tags, strings.TrimSpace(t))
			}
		}
	}

	input := NoteInput{
		ProjectID: projectID,
		Title:     r.FormValue("title"),
		Content:   r.FormValue("content"),
		IsPinned:  r.FormValue("is_pinned") == "true",
		Tags:      tags,
	}

	n, err := h.service.CreateNote(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to create note", zap.Error(err))
		http.Redirect(w, r, "/notes?error=failed_to_create", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/notes/"+strconv.FormatInt(n.ID, 10), http.StatusSeeOther)
}

// GET /notes/{id}
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	noteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	noteDetail, err := h.service.GetNote(r.Context(), user.ID, noteID)
	if err != nil {
		h.logger.Error("note not found", zap.Error(err))
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
		return
	}

	var availableProjects []db.Project
	if h.projectsRepo != nil {
		availableProjects, _ = h.projectsRepo.ListProjects(r.Context(), user.ID)
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return notetemplates.Editor(*noteDetail, availableProjects, false).Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base(noteDetail.Note.Title, "/notes", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /notes/{id}/edit
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	noteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}

	var projectID int64
	if projStr := r.FormValue("project_id"); projStr != "" {
		projectID, _ = strconv.ParseInt(projStr, 10, 64)
	}

	tagsRaw := r.FormValue("tags")
	var tags []string
	if tagsRaw != "" {
		for _, t := range strings.Split(tagsRaw, ",") {
			if strings.TrimSpace(t) != "" {
				tags = append(tags, strings.TrimSpace(t))
			}
		}
	}

	input := NoteInput{
		ProjectID: projectID,
		Title:     r.FormValue("title"),
		Content:   r.FormValue("content"),
		IsPinned:  r.FormValue("is_pinned") == "true",
		Tags:      tags,
	}

	_, err = h.service.UpdateNote(r.Context(), user.ID, noteID, input)
	if err != nil {
		h.logger.Error("failed to update note", zap.Error(err))
		http.Redirect(w, r, "/notes/"+strconv.FormatInt(noteID, 10)+"?error=update_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/notes/"+strconv.FormatInt(noteID, 10), http.StatusSeeOther)
}

// POST /notes/{id}/autosave
func (h *Handler) Autosave(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	noteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := NoteInput{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
	}

	_, err = h.service.UpdateNote(r.Context(), user.ID, noteID, input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<span class="text-xs text-emerald-400 font-mono-code">Autosaved</span>`))
}

// POST /notes/{id}/pin
func (h *Handler) TogglePin(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	noteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	_, err = h.service.TogglePin(r.Context(), user.ID, noteID)
	if err != nil {
		h.logger.Error("failed to toggle note pin", zap.Error(err))
	}

	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}

// POST /notes/{id}/delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	noteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteNote(r.Context(), user.ID, noteID)
	if err != nil {
		h.logger.Error("failed to delete note", zap.Error(err))
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/notes", http.StatusSeeOther)
}
