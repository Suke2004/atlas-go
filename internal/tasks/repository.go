package tasks

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
	GetTaskByID(ctx context.Context, arg db.GetTaskByIDParams) (db.Task, error)
	ListTasks(ctx context.Context, userID int64) ([]db.Task, error)
	ListTasksByProject(ctx context.Context, arg db.ListTasksByProjectParams) ([]db.Task, error)
	UpdateTaskStatus(ctx context.Context, arg db.UpdateTaskStatusParams) (db.Task, error)
	UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error)
	DeleteTask(ctx context.Context, arg db.DeleteTaskParams) error

	AddTaskLabel(ctx context.Context, arg db.AddTaskLabelParams) error
	ListTaskLabels(ctx context.Context, taskID int64) ([]string, error)

	AddTaskDependency(ctx context.Context, arg db.AddTaskDependencyParams) error
	ListTaskDependencies(ctx context.Context, taskID int64) ([]db.Task, error)
	GetTodayFocusTasks(ctx context.Context, userID int64) ([]db.Task, error)
}

type repository struct {
	database *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return r.database.Queries.CreateTask(ctx, arg)
}

func (r *repository) GetTaskByID(ctx context.Context, arg db.GetTaskByIDParams) (db.Task, error) {
	return r.database.Queries.GetTaskByID(ctx, arg)
}

func (r *repository) ListTasks(ctx context.Context, userID int64) ([]db.Task, error) {
	return r.database.Queries.ListTasks(ctx, userID)
}

func (r *repository) ListTasksByProject(ctx context.Context, arg db.ListTasksByProjectParams) ([]db.Task, error) {
	return r.database.Queries.ListTasksByProject(ctx, arg)
}

func (r *repository) UpdateTaskStatus(ctx context.Context, arg db.UpdateTaskStatusParams) (db.Task, error) {
	return r.database.Queries.UpdateTaskStatus(ctx, arg)
}

func (r *repository) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	return r.database.Queries.UpdateTask(ctx, arg)
}

func (r *repository) DeleteTask(ctx context.Context, arg db.DeleteTaskParams) error {
	return r.database.Queries.DeleteTask(ctx, arg)
}

func (r *repository) AddTaskLabel(ctx context.Context, arg db.AddTaskLabelParams) error {
	return r.database.Queries.AddTaskLabel(ctx, arg)
}

func (r *repository) ListTaskLabels(ctx context.Context, taskID int64) ([]string, error) {
	return r.database.Queries.ListTaskLabels(ctx, taskID)
}

func (r *repository) AddTaskDependency(ctx context.Context, arg db.AddTaskDependencyParams) error {
	return r.database.Queries.AddTaskDependency(ctx, arg)
}

func (r *repository) ListTaskDependencies(ctx context.Context, taskID int64) ([]db.Task, error) {
	return r.database.Queries.ListTaskDependencies(ctx, taskID)
}

func (r *repository) GetTodayFocusTasks(ctx context.Context, userID int64) ([]db.Task, error) {
	return r.database.Queries.GetTodayFocusTasks(ctx, userID)
}
