package db

import (
	"context"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (
    user_id, name, description, status, color, target_date,
    github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, user_id, name, description, status, color, progress_percentage, target_date, github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack, created_at, updated_at
`

type CreateProjectParams struct {
	UserID             int64       `json:"user_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Status             string      `json:"status"`
	Color              string      `json:"color"`
	TargetDate         interface{} `json:"target_date"`
	GithubUrl          string      `json:"github_url"`
	GithubStars        int64       `json:"github_stars"`
	GithubForks        int64       `json:"github_forks"`
	GithubOpenIssues   int64       `json:"github_open_issues"`
	GithubLanguage     string      `json:"github_language"`
	GithubLastPushedAt interface{} `json:"github_last_pushed_at"`
	TechStack          string      `json:"tech_stack"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject,
		arg.UserID,
		arg.Name,
		arg.Description,
		arg.Status,
		arg.Color,
		arg.TargetDate,
		arg.GithubUrl,
		arg.GithubStars,
		arg.GithubForks,
		arg.GithubOpenIssues,
		arg.GithubLanguage,
		arg.GithubLastPushedAt,
		arg.TechStack,
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
		&i.GithubUrl,
		&i.GithubStars,
		&i.GithubForks,
		&i.GithubOpenIssues,
		&i.GithubLanguage,
		&i.GithubLastPushedAt,
		&i.TechStack,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProjectByID = `-- name: GetProjectByID :one
SELECT id, user_id, name, description, status, color, progress_percentage, target_date, github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack, created_at, updated_at FROM projects
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
		&i.GithubUrl,
		&i.GithubStars,
		&i.GithubForks,
		&i.GithubOpenIssues,
		&i.GithubLanguage,
		&i.GithubLastPushedAt,
		&i.TechStack,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProjects = `-- name: ListProjects :many
SELECT id, user_id, name, description, status, color, progress_percentage, target_date, github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack, created_at, updated_at FROM projects
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
			&i.GithubUrl,
			&i.GithubStars,
			&i.GithubForks,
			&i.GithubOpenIssues,
			&i.GithubLanguage,
			&i.GithubLastPushedAt,
			&i.TechStack,
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

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET name = ?, description = ?, status = ?, color = ?, target_date = ?, github_url = ?, tech_stack = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING id, user_id, name, description, status, color, progress_percentage, target_date, github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack, created_at, updated_at
`

type UpdateProjectParams struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Color       string      `json:"color"`
	TargetDate  interface{} `json:"target_date"`
	GithubUrl   string      `json:"github_url"`
	TechStack   string      `json:"tech_stack"`
	ID          int64       `json:"id"`
	UserID      int64       `json:"user_id"`
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, updateProject,
		arg.Name,
		arg.Description,
		arg.Status,
		arg.Color,
		arg.TargetDate,
		arg.GithubUrl,
		arg.TechStack,
		arg.ID,
		arg.UserID,
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
		&i.GithubUrl,
		&i.GithubStars,
		&i.GithubForks,
		&i.GithubOpenIssues,
		&i.GithubLanguage,
		&i.GithubLastPushedAt,
		&i.TechStack,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateGitHubStats = `-- name: UpdateGitHubStats :one
UPDATE projects
SET github_stars = ?, github_forks = ?, github_open_issues = ?, github_language = ?, github_last_pushed_at = ?, tech_stack = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING id, user_id, name, description, status, color, progress_percentage, target_date, github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack, created_at, updated_at
`

type UpdateGitHubStatsParams struct {
	GithubStars        int64       `json:"github_stars"`
	GithubForks        int64       `json:"github_forks"`
	GithubOpenIssues   int64       `json:"github_open_issues"`
	GithubLanguage     string      `json:"github_language"`
	GithubLastPushedAt interface{} `json:"github_last_pushed_at"`
	TechStack          string      `json:"tech_stack"`
	ID                 int64       `json:"id"`
	UserID             int64       `json:"user_id"`
}

func (q *Queries) UpdateGitHubStats(ctx context.Context, arg UpdateGitHubStatsParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, updateGitHubStats,
		arg.GithubStars,
		arg.GithubForks,
		arg.GithubOpenIssues,
		arg.GithubLanguage,
		arg.GithubLastPushedAt,
		arg.TechStack,
		arg.ID,
		arg.UserID,
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
		&i.GithubUrl,
		&i.GithubStars,
		&i.GithubForks,
		&i.GithubOpenIssues,
		&i.GithubLanguage,
		&i.GithubLastPushedAt,
		&i.TechStack,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
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

const deleteProject = `-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ? AND user_id = ?
`

type DeleteProjectParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) DeleteProject(ctx context.Context, arg DeleteProjectParams) error {
	_, err := q.db.ExecContext(ctx, deleteProject, arg.ID, arg.UserID)
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

const listMilestonesByProject = `-- name: ListMilestonesByProject :many
SELECT id, project_id, title, due_date, is_completed, completed_at, created_at FROM milestones
WHERE project_id = ?
ORDER BY created_at ASC
`

func (q *Queries) ListMilestonesByProject(ctx context.Context, projectID int64) ([]Milestone, error) {
	rows, err := q.db.QueryContext(ctx, listMilestonesByProject, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Milestone
	for rows.Next() {
		var i Milestone
		if err := rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.Title,
			&i.DueDate,
			&i.IsCompleted,
			&i.CompletedAt,
			&i.CreatedAt,
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

const toggleMilestone = `-- name: ToggleMilestone :one
UPDATE milestones
SET is_completed = ?, completed_at = CASE WHEN ? = 1 THEN CURRENT_TIMESTAMP ELSE NULL END
WHERE id = ? AND project_id = ?
RETURNING id, project_id, title, due_date, is_completed, completed_at, created_at
`

type ToggleMilestoneParams struct {
	IsCompleted bool  `json:"is_completed"`
	IsCompleted_2 bool `json:"is_completed_2"`
	ID          int64 `json:"id"`
	ProjectID   int64 `json:"project_id"`
}

func (q *Queries) ToggleMilestone(ctx context.Context, arg ToggleMilestoneParams) (Milestone, error) {
	row := q.db.QueryRowContext(ctx, toggleMilestone, arg.IsCompleted, arg.IsCompleted_2, arg.ID, arg.ProjectID)
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

const deleteMilestone = `-- name: DeleteMilestone :exec
DELETE FROM milestones
WHERE id = ? AND project_id = ?
`

type DeleteMilestoneParams struct {
	ID        int64 `json:"id"`
	ProjectID int64 `json:"project_id"`
}

func (q *Queries) DeleteMilestone(ctx context.Context, arg DeleteMilestoneParams) error {
	_, err := q.db.ExecContext(ctx, deleteMilestone, arg.ID, arg.ProjectID)
	return err
}
