package db

import (
	"context"
)

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes, actual_minutes, completed_at, created_at, updated_at
`

type CreateTaskParams struct {
	UserID           int64       `json:"user_id"`
	ProjectID        interface{} `json:"project_id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Status           string      `json:"status"`
	Priority         string      `json:"priority"`
	EnergyLevel      string      `json:"energy_level"`
	DueDate          interface{} `json:"due_date"`
	EstimatedMinutes interface{} `json:"estimated_minutes"`
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, createTask,
		arg.UserID,
		arg.ProjectID,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.Priority,
		arg.EnergyLevel,
		arg.DueDate,
		arg.EstimatedMinutes,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProjectID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.Priority,
		&i.EnergyLevel,
		&i.DueDate,
		&i.EstimatedMinutes,
		&i.ActualMinutes,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTaskByID = `-- name: GetTaskByID :one
SELECT id, user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes, actual_minutes, completed_at, created_at, updated_at FROM tasks
WHERE id = ? AND user_id = ? LIMIT 1
`

type GetTaskByIDParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetTaskByID(ctx context.Context, arg GetTaskByIDParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskByID, arg.ID, arg.UserID)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProjectID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.Priority,
		&i.EnergyLevel,
		&i.DueDate,
		&i.EstimatedMinutes,
		&i.ActualMinutes,
		&i.CompletedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTasks = `-- name: ListTasks :many
SELECT id, user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes, actual_minutes, completed_at, created_at, updated_at FROM tasks
WHERE user_id = ?
ORDER BY 
  CASE status WHEN 'todo' THEN 1 WHEN 'in_progress' THEN 2 WHEN 'done' THEN 3 END,
  CASE priority WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END,
  created_at DESC
`

func (q *Queries) ListTasks(ctx context.Context, userID int64) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProjectID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.Priority,
			&i.EnergyLevel,
			&i.DueDate,
			&i.EstimatedMinutes,
			&i.ActualMinutes,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const addTaskLabel = `-- name: AddTaskLabel :exec
INSERT INTO task_labels (task_id, label)
VALUES (?, ?)
ON CONFLICT (task_id, label) DO NOTHING
`

type AddTaskLabelParams struct {
	TaskID int64  `json:"task_id"`
	Label  string `json:"label"`
}

func (q *Queries) AddTaskLabel(ctx context.Context, arg AddTaskLabelParams) error {
	_, err := q.db.ExecContext(ctx, addTaskLabel, arg.TaskID, arg.Label)
	return err
}

const getTodayFocusTasks = `-- name: GetTodayFocusTasks :many
SELECT id, user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes, actual_minutes, completed_at, created_at, updated_at FROM tasks
WHERE user_id = ? AND status != 'done' AND (date(due_date) = date('now') OR priority IN ('critical', 'high'))
ORDER BY priority DESC, due_date ASC
LIMIT 3
`

func (q *Queries) GetTodayFocusTasks(ctx context.Context, userID int64) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTodayFocusTasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProjectID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.Priority,
			&i.EnergyLevel,
			&i.DueDate,
			&i.EstimatedMinutes,
			&i.ActualMinutes,
			&i.CompletedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
