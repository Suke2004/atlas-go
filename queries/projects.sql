-- name: CreateProject :one
INSERT INTO projects (
    user_id, name, description, status, color, target_date,
    github_url, github_stars, github_forks, github_open_issues, github_language, github_last_pushed_at, tech_stack
) VALUES (
    ?, ?, ?, ?, ?, ?,
    ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = ? AND user_id = ? LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: ListProjectsByStatus :many
SELECT * FROM projects
WHERE user_id = ? AND status = ?
ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = ?, description = ?, status = ?, color = ?, target_date = ?,
    github_url = ?, tech_stack = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateGitHubStats :one
UPDATE projects
SET github_stars = ?, github_forks = ?, github_open_issues = ?, github_language = ?, github_last_pushed_at = ?, tech_stack = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateProjectProgress :exec
UPDATE projects
SET progress_percentage = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = ? AND user_id = ?;

-- name: CreateMilestone :one
INSERT INTO milestones (project_id, title, due_date)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListMilestonesByProject :many
SELECT * FROM milestones
WHERE project_id = ?
ORDER BY created_at ASC;

-- name: ToggleMilestone :one
UPDATE milestones
SET is_completed = ?, completed_at = CASE WHEN ? = 1 THEN CURRENT_TIMESTAMP ELSE NULL END
WHERE id = ? AND project_id = ?
RETURNING *;

-- name: DeleteMilestone :exec
DELETE FROM milestones
WHERE id = ? AND project_id = ?;
