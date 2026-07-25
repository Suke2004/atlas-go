-- +goose Up
-- Migration: 004_tasks.sql
-- Purpose: Tasks management including task priority, energy level, estimate vs actual,
-- labels (tags), and task dependencies.

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo', -- 'todo' | 'in_progress' | 'done'
    priority TEXT NOT NULL DEFAULT 'medium', -- 'low' | 'medium' | 'high' | 'critical'
    energy_level TEXT NOT NULL DEFAULT 'medium', -- 'low' | 'medium' | 'high'
    due_date DATETIME,
    estimated_minutes INTEGER DEFAULT 0,
    actual_minutes INTEGER DEFAULT 0,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS task_labels (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    PRIMARY KEY (task_id, label)
);

CREATE TABLE IF NOT EXISTS task_dependencies (
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, depends_on_id)
);

CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);

-- +goose Down
DROP TABLE IF EXISTS task_dependencies;
DROP TABLE IF EXISTS task_labels;
DROP TABLE IF EXISTS tasks;
