package db

import (
	"context"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (user_id, name, description, status, color, target_date)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, user_id, name, description, status, color, progress_percentage, target_date, created_at, updated_at
`

type CreateProjectParams struct {
	UserID      int64       `json:"user_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Color       string      `json:"color"`
	TargetDate  interface{} `json:"target_date"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.UserID,
		arg.Name,
		arg.Description,
		arg.Status,
		arg.Color,
		arg.TargetDate,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Status,
		&i.Color,
		&i.ProgressPercentage,
		&i.TargetDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProjectByID = `-- name: GetProjectByID :one
SELECT id, user_id, name, description, status, color, progress_percentage, target_date, created_at, updated_at FROM projects
WHERE id = ? AND user_id = ? LIMIT 1
`

type GetProjectByIDParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetProjectByID(ctx context.Context, arg GetProjectByIDParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, getProjectByID, arg.ID, arg.UserID)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Status,
		&i.Color,
		&i.ProgressPercentage,
		&i.TargetDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProjects = `-- name: ListProjects :many
SELECT id, user_id, name, description, status, color, progress_percentage, target_date, created_at, updated_at FROM projects
WHERE user_id = ?
ORDER BY created_at DESC
`

func (q *Queries) ListProjects(ctx context.Context, userID int64) ([]Project, error) {
	rows, err := q.db.QueryContext(ctx, listProjects, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.Status,
			&i.Color,
			&i.ProgressPercentage,
			&i.TargetDate,
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

const updateProjectProgress = `-- name: UpdateProjectProgress :exec
UPDATE projects
SET progress_percentage = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
`

type UpdateProjectProgressParams struct {
	ProgressPercentage int64 `json:"progress_percentage"`
	ID                 int64 `json:"id"`
	UserID             int64 `json:"user_id"`
}

func (q *Queries) UpdateProjectProgress(ctx context.Context, arg UpdateProjectProgressParams) error {
	_, err := q.db.ExecContext(ctx, updateProjectProgress, arg.ProgressPercentage, arg.ID, arg.UserID)
	return err
}

const createMilestone = `-- name: CreateMilestone :one
INSERT INTO milestones (project_id, title, due_date)
VALUES (?, ?, ?)
RETURNING id, project_id, title, due_date, is_completed, completed_at, created_at
`

type CreateMilestoneParams struct {
	ProjectID int64       `json:"project_id"`
	Title     string      `json:"title"`
	DueDate   interface{} `json:"due_date"`
}

func (q *Queries) CreateMilestone(ctx context.Context, arg CreateMilestoneParams) (Milestone, error) {
	row := q.db.QueryRowContext(ctx, createMilestone, arg.ProjectID, arg.Title, arg.DueDate)
	var i Milestone
	err := row.Scan(
		&i.ID,
		&i.ProjectID,
		&i.Title,
		&i.DueDate,
		&i.IsCompleted,
		&i.CompletedAt,
		&i.CreatedAt,
	)
	return i, err
}
