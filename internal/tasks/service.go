package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	tasktemplates "github.com/Suke2004/atlas-go/web/templates/tasks"
)

type TaskInput struct {
	ProjectID        int64    `json:"project_id"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Status           string   `json:"status"`
	Priority         string   `json:"priority"`
	EnergyLevel      string   `json:"energy_level"`
	DueDate          string   `json:"due_date"`
	EstimatedMinutes int64    `json:"estimated_minutes"`
	ActualMinutes    int64    `json:"actual_minutes"`
	Labels           []string `json:"labels"`
}

type TaskWithDetails = tasktemplates.TaskWithDetails
type TasksSummary = tasktemplates.TasksSummary

type Service interface {
	CreateTask(ctx context.Context, userID int64, input TaskInput) (db.Task, error)
	GetTask(ctx context.Context, userID, taskID int64) (*TaskWithDetails, error)
	ListTasks(ctx context.Context, userID int64, statusFilter, priorityFilter, energyFilter string, projectID int64, searchQuery string) ([]TaskWithDetails, error)
	GetTasksSummary(ctx context.Context, userID int64) (TasksSummary, error)
	UpdateTaskStatus(ctx context.Context, userID, taskID int64, status string) (db.Task, error)
	UpdateTask(ctx context.Context, userID, taskID int64, input TaskInput) (db.Task, error)
	DeleteTask(ctx context.Context, userID, taskID int64) error
}

type service struct {
	repo           Repository
	projectsRepo   projects.Repository
	projectsSvc    projects.Service
}

func NewService(repo Repository, projRepo projects.Repository, projSvc projects.Service) Service {
	return &service{
		repo:         repo,
		projectsRepo: projRepo,
		projectsSvc:  projSvc,
	}
}

func (s *service) CreateTask(ctx context.Context, userID int64, input TaskInput) (db.Task, error) {
	if strings.TrimSpace(input.Title) == "" {
		return db.Task{}, fmt.Errorf("task title is required")
	}

	if input.Status == "" {
		input.Status = "todo"
	}
	if input.Priority == "" {
		input.Priority = "medium"
	}
	if input.EnergyLevel == "" {
		input.EnergyLevel = "medium"
	}

	var projectID sql.NullInt64
	if input.ProjectID > 0 {
		projectID = sql.NullInt64{Int64: input.ProjectID, Valid: true}
	}

	var dueDate interface{}
	if strings.TrimSpace(input.DueDate) != "" {
		dueDate = strings.TrimSpace(input.DueDate)
	}

	var estMinutes sql.NullInt64
	if input.EstimatedMinutes > 0 {
		estMinutes = sql.NullInt64{Int64: input.EstimatedMinutes, Valid: true}
	}

	arg := db.CreateTaskParams{
		UserID:           userID,
		ProjectID:        projectID,
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		Status:           input.Status,
		Priority:         input.Priority,
		EnergyLevel:      input.EnergyLevel,
		DueDate:          dueDate,
		EstimatedMinutes: estMinutes,
	}

	t, err := s.repo.CreateTask(ctx, arg)
	if err != nil {
		return db.Task{}, err
	}

	// Add labels
	for _, label := range input.Labels {
		lbl := strings.TrimSpace(label)
		if lbl != "" {
			_ = s.repo.AddTaskLabel(ctx, db.AddTaskLabelParams{
				TaskID: t.ID,
				Label:  lbl,
			})
		}
	}

	// Recalculate project progress if linked to project
	if input.ProjectID > 0 && s.projectsSvc != nil {
		_ = s.projectsSvc.RecalculateProgress(ctx, userID, input.ProjectID)
	}

	return t, nil
}

func (s *service) GetTask(ctx context.Context, userID, taskID int64) (*TaskWithDetails, error) {
	t, err := s.repo.GetTaskByID(ctx, db.GetTaskByIDParams{ID: taskID, UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	var projName string
	if t.ProjectID.Valid && s.projectsRepo != nil {
		proj, err := s.projectsRepo.GetProjectByID(ctx, db.GetProjectByIDParams{ID: t.ProjectID.Int64, UserID: userID})
		if err == nil {
			projName = proj.Name
		}
	}

	labels, _ := s.repo.ListTaskLabels(ctx, taskID)
	dependencies, _ := s.repo.ListTaskDependencies(ctx, taskID)

	return &TaskWithDetails{
		Task:         t,
		ProjectName:  projName,
		Labels:       labels,
		Dependencies: dependencies,
	}, nil
}

func (s *service) ListTasks(ctx context.Context, userID int64, statusFilter, priorityFilter, energyFilter string, projectID int64, searchQuery string) ([]TaskWithDetails, error) {
	allTasks, err := s.repo.ListTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	var projectsMap map[int64]string
	if s.projectsRepo != nil {
		projList, err := s.projectsRepo.ListProjects(ctx, userID)
		if err == nil {
			projectsMap = make(map[int64]string)
			for _, p := range projList {
				projectsMap[p.ID] = p.Name
			}
		}
	}

	var filtered []TaskWithDetails
	searchLower := strings.ToLower(strings.TrimSpace(searchQuery))

	for _, t := range allTasks {
		if statusFilter != "" && statusFilter != "all" && !strings.EqualFold(t.Status, statusFilter) {
			continue
		}
		if priorityFilter != "" && priorityFilter != "all" && !strings.EqualFold(t.Priority, priorityFilter) {
			continue
		}
		if energyFilter != "" && energyFilter != "all" && !strings.EqualFold(t.EnergyLevel, energyFilter) {
			continue
		}
		if projectID > 0 && (!t.ProjectID.Valid || t.ProjectID.Int64 != projectID) {
			continue
		}
		if searchLower != "" {
			titleMatch := strings.Contains(strings.ToLower(t.Title), searchLower)
			descMatch := strings.Contains(strings.ToLower(t.Description), searchLower)
			if !titleMatch && !descMatch {
				continue
			}
		}

		projName := ""
		if t.ProjectID.Valid && projectsMap != nil {
			projName = projectsMap[t.ProjectID.Int64]
		}

		labels, _ := s.repo.ListTaskLabels(ctx, t.ID)

		filtered = append(filtered, TaskWithDetails{
			Task:        t,
			ProjectName: projName,
			Labels:      labels,
		})
	}

	return filtered, nil
}

func (s *service) GetTasksSummary(ctx context.Context, userID int64) (TasksSummary, error) {
	allTasks, err := s.repo.ListTasks(ctx, userID)
	if err != nil {
		return TasksSummary{}, err
	}

	var summary TasksSummary
	summary.TotalTasks = int64(len(allTasks))

	for _, t := range allTasks {
		switch strings.ToLower(t.Status) {
		case "todo":
			summary.TodoTasks++
		case "in_progress":
			summary.InProgressTasks++
		case "done":
			summary.DoneTasks++
		}
	}

	focus, err := s.repo.GetTodayFocusTasks(ctx, userID)
	if err == nil {
		summary.TodayFocusTasks = int64(len(focus))
	}

	if summary.TotalTasks > 0 {
		summary.CompletionRate = (summary.DoneTasks * 100) / summary.TotalTasks
	}

	return summary, nil
}

func (s *service) UpdateTaskStatus(ctx context.Context, userID, taskID int64, status string) (db.Task, error) {
	t, err := s.repo.UpdateTaskStatus(ctx, db.UpdateTaskStatusParams{
		Status: status,
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		return db.Task{}, err
	}

	if t.ProjectID.Valid && s.projectsSvc != nil {
		_ = s.projectsSvc.RecalculateProgress(ctx, userID, t.ProjectID.Int64)
	}

	return t, nil
}

func (s *service) UpdateTask(ctx context.Context, userID, taskID int64, input TaskInput) (db.Task, error) {
	if strings.TrimSpace(input.Title) == "" {
		return db.Task{}, fmt.Errorf("task title is required")
	}

	var projectID sql.NullInt64
	if input.ProjectID > 0 {
		projectID = sql.NullInt64{Int64: input.ProjectID, Valid: true}
	}

	var dueDate interface{}
	if strings.TrimSpace(input.DueDate) != "" {
		dueDate = strings.TrimSpace(input.DueDate)
	}

	var estMinutes, actMinutes sql.NullInt64
	if input.EstimatedMinutes > 0 {
		estMinutes = sql.NullInt64{Int64: input.EstimatedMinutes, Valid: true}
	}
	if input.ActualMinutes > 0 {
		actMinutes = sql.NullInt64{Int64: input.ActualMinutes, Valid: true}
	}

	t, err := s.repo.UpdateTask(ctx, db.UpdateTaskParams{
		ProjectID:        projectID,
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		Status:           input.Status,
		Priority:         input.Priority,
		EnergyLevel:      input.EnergyLevel,
		DueDate:          dueDate,
		EstimatedMinutes: estMinutes,
		ActualMinutes:    actMinutes,
		ID:               taskID,
		UserID:           userID,
	})
	if err != nil {
		return db.Task{}, err
	}

	if input.ProjectID > 0 && s.projectsSvc != nil {
		_ = s.projectsSvc.RecalculateProgress(ctx, userID, input.ProjectID)
	}

	return t, nil
}

func (s *service) DeleteTask(ctx context.Context, userID, taskID int64) error {
	t, err := s.repo.GetTaskByID(ctx, db.GetTaskByIDParams{ID: taskID, UserID: userID})
	if err != nil {
		return err
	}

	err = s.repo.DeleteTask(ctx, db.DeleteTaskParams{ID: taskID, UserID: userID})
	if err != nil {
		return err
	}

	if t.ProjectID.Valid && s.projectsSvc != nil {
		_ = s.projectsSvc.RecalculateProgress(ctx, userID, t.ProjectID.Int64)
	}

	return nil
}
