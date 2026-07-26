package projects

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
	projecttemplates "github.com/Suke2004/atlas-go/web/templates/projects"
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

// GET /projects
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	statusFilter := r.URL.Query().Get("status")
	if statusFilter == "" {
		statusFilter = "all"
	}
	tagFilter := r.URL.Query().Get("tag")
	searchQuery := r.URL.Query().Get("search")
	viewMode := r.URL.Query().Get("view")
	if viewMode == "" {
		viewMode = "grid"
	}

	projectList, err := h.service.ListProjects(r.Context(), user.ID, statusFilter, tagFilter, searchQuery)
	if err != nil {
		h.logger.Error("failed to list projects", zap.Error(err))
		http.Error(w, "Failed to load projects", http.StatusInternalServerError)
		return
	}

	summary, err := h.service.GetProjectsSummary(r.Context(), user.ID)
	if err != nil {
		h.logger.Warn("failed to calculate projects summary", zap.Error(err))
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX && r.Header.Get("HX-Target") == "projects-content" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = projecttemplates.ProjectContent(projectList, statusFilter, viewMode, tagFilter, searchQuery).Render(r.Context(), w)
		return
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return projecttemplates.List(projectList, summary, statusFilter, viewMode, tagFilter, searchQuery, "").Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Projects", "/projects", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /projects
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	input := ProjectInput{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
		Color:       r.FormValue("color"),
		TargetDate:  r.FormValue("target_date"),
		GithubURL:   r.FormValue("github_url"),
		TechStack:   r.FormValue("tech_stack"),
	}

	_, err := h.service.CreateProject(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to create project", zap.Error(err))
		http.Redirect(w, r, "/projects?error=failed_to_create", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

// GET /projects/{id}
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	data, err := h.service.GetProject(r.Context(), user.ID, projectID)
	if err != nil {
		h.logger.Error("failed to get project detail", zap.Error(err))
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return projecttemplates.Detail(data.Project, data.Milestones, "").Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base(data.Project.Name, "/projects", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /projects/{id}/edit
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	input := ProjectInput{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
		Color:       r.FormValue("color"),
		TargetDate:  r.FormValue("target_date"),
		GithubURL:   r.FormValue("github_url"),
		TechStack:   r.FormValue("tech_stack"),
	}

	_, err = h.service.UpdateProject(r.Context(), user.ID, projectID, input)
	if err != nil {
		h.logger.Error("failed to update project", zap.Error(err))
	}

	http.Redirect(w, r, "/projects/"+strconv.FormatInt(projectID, 10), http.StatusSeeOther)
}

// POST /projects/{id}/delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	_ = h.service.DeleteProject(r.Context(), user.ID, projectID)
	http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

// POST /projects/{id}/sync-github
func (h *Handler) SyncGitHub(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	updatedProject, err := h.service.SyncGitHubStats(r.Context(), user.ID, projectID)
	if err != nil {
		h.logger.Error("failed to sync github stats", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = projecttemplates.GitHubCard(updatedProject).Render(r.Context(), w)
}

// POST /projects/{id}/import-issues
func (h *Handler) ImportIssues(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	_, err = h.service.ImportGitHubIssues(r.Context(), user.ID, projectID)
	if err != nil {
		h.logger.Error("failed to import github issues", zap.Error(err))
	}

	http.Redirect(w, r, "/projects/"+strconv.FormatInt(projectID, 10), http.StatusSeeOther)
}

// POST /projects/{id}/milestones
func (h *Handler) CreateMilestone(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	dueDate := r.FormValue("due_date")

	_, err = h.service.CreateMilestone(r.Context(), user.ID, projectID, title, dueDate)
	if err != nil {
		h.logger.Error("failed to create milestone", zap.Error(err))
	}

	http.Redirect(w, r, "/projects/"+strconv.FormatInt(projectID, 10), http.StatusSeeOther)
}

// POST /projects/{id}/milestones/{milestoneID}/toggle
func (h *Handler) ToggleMilestone(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	milestoneID, err := strconv.ParseInt(chi.URLParam(r, "milestoneID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid milestone ID", http.StatusBadRequest)
		return
	}

	_ = r.ParseForm()
	isCompleted := r.FormValue("completed") == "true" || r.FormValue("completed") == "on" || r.FormValue("completed") == "1"

	_, err = h.service.ToggleMilestone(r.Context(), user.ID, projectID, milestoneID, isCompleted)
	if err != nil {
		h.logger.Error("failed to toggle milestone", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/projects/"+strconv.FormatInt(projectID, 10), http.StatusSeeOther)
}

// POST /projects/{id}/milestones/{milestoneID}/delete
func (h *Handler) DeleteMilestone(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	milestoneID, err := strconv.ParseInt(chi.URLParam(r, "milestoneID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid milestone ID", http.StatusBadRequest)
		return
	}

	_ = h.service.DeleteMilestone(r.Context(), user.ID, projectID, milestoneID)
	http.Redirect(w, r, "/projects/"+strconv.FormatInt(projectID, 10), http.StatusSeeOther)
}
