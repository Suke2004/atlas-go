-- name: CreateTask :one
INSERT INTO tasks (user_id, project_id, title, description, status, priority, energy_level, due_date, estimated_minutes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = ? AND user_id = ? LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
WHERE user_id = ?
ORDER BY 
  CASE status WHEN 'todo' THEN 1 WHEN 'in_progress' THEN 2 WHEN 'done' THEN 3 END,
  CASE priority WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END,
  created_at DESC;

-- name: ListTasksByProject :many
SELECT * FROM tasks
WHERE user_id = ? AND project_id = ?
ORDER BY created_at DESC;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET status = ?, completed_at = CASE WHEN ? = 'done' THEN CURRENT_TIMESTAMP ELSE NULL END, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET project_id = ?, title = ?, description = ?, status = ?, priority = ?, energy_level = ?, due_date = ?, estimated_minutes = ?, actual_minutes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ? AND user_id = ?;

-- name: AddTaskLabel :exec
INSERT INTO task_labels (task_id, label)
VALUES (?, ?)
ON CONFLICT (task_id, label) DO NOTHING;

-- name: ListTaskLabels :many
SELECT label FROM task_labels
WHERE task_id = ?;

-- name: AddTaskDependency :exec
INSERT INTO task_dependencies (task_id, depends_on_id)
VALUES (?, ?)
ON CONFLICT (task_id, depends_on_id) DO NOTHING;

-- name: ListTaskDependencies :many
SELECT t.* FROM tasks t
JOIN task_dependencies td ON t.id = td.depends_on_id
WHERE td.task_id = ?;

-- name: GetTodayFocusTasks :many
SELECT * FROM tasks
WHERE user_id = ? AND status != 'done' AND (date(due_date) = date('now') OR priority IN ('critical', 'high'))
ORDER BY priority DESC, due_date ASC
LIMIT 3;
