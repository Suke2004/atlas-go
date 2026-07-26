package tasks

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/Suke2004/atlas-go/internal/auth"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	layouts "github.com/Suke2004/atlas-go/web/templates/layout"
	tasktemplates "github.com/Suke2004/atlas-go/web/templates/tasks"
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

// GET /tasks
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	statusFilter := r.URL.Query().Get("status")
	priorityFilter := r.URL.Query().Get("priority")
	energyFilter := r.URL.Query().Get("energy")
	searchQuery := r.URL.Query().Get("search")
	viewMode := r.URL.Query().Get("view")
	if viewMode == "" {
		viewMode = "list"
	}

	var projectID int64
	if projStr := r.URL.Query().Get("project_id"); projStr != "" {
		projectID, _ = strconv.ParseInt(projStr, 10, 64)
	}

	taskList, err := h.service.ListTasks(r.Context(), user.ID, statusFilter, priorityFilter, energyFilter, projectID, searchQuery)
	if err != nil {
		h.logger.Error("failed to list tasks", zap.Error(err))
		http.Error(w, "Failed to load tasks", http.StatusInternalServerError)
		return
	}

	summary, err := h.service.GetTasksSummary(r.Context(), user.ID)
	if err != nil {
		h.logger.Warn("failed to calculate tasks summary", zap.Error(err))
	}

	var availableProjects []db.Project
	if h.projectsRepo != nil {
		availableProjects, _ = h.projectsRepo.ListProjects(r.Context(), user.ID)
	}

	isHTMX := r.Header.Get("HX-Request") == "true"
	if isHTMX && r.Header.Get("HX-Target") == "tasks-content" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tasktemplates.TaskContent(taskList, availableProjects, statusFilter, viewMode).Render(r.Context(), w)
		return
	}

	pageContent := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return tasktemplates.List(taskList, summary, availableProjects, statusFilter, priorityFilter, energyFilter, viewMode, searchQuery, "").Render(ctx, w)
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = layouts.Base("Tasks", "/tasks", user.DisplayName, pageContent).Render(r.Context(), w)
}

// POST /tasks
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

	var estMins int64
	if estStr := r.FormValue("estimated_minutes"); estStr != "" {
		estMins, _ = strconv.ParseInt(estStr, 10, 64)
	}

	input := TaskInput{
		ProjectID:        projectID,
		Title:            r.FormValue("title"),
		Description:      r.FormValue("description"),
		Status:           r.FormValue("status"),
		Priority:         r.FormValue("priority"),
		EnergyLevel:      r.FormValue("energy_level"),
		DueDate:          r.FormValue("due_date"),
		EstimatedMinutes: estMins,
	}

	_, err := h.service.CreateTask(r.Context(), user.ID, input)
	if err != nil {
		h.logger.Error("failed to create task", zap.Error(err))
		http.Redirect(w, r, "/tasks?error=failed_to_create", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

// POST /tasks/{id}/status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if status == "" {
		status = "done"
	}

	_, err = h.service.UpdateTaskStatus(r.Context(), user.ID, taskID, status)
	if err != nil {
		h.logger.Error("failed to update task status", zap.Error(err))
		http.Error(w, "Failed to update task status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

// POST /tasks/{id}/delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	taskID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteTask(r.Context(), user.ID, taskID)
	if err != nil {
		h.logger.Error("failed to delete task", zap.Error(err))
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}
