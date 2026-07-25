-- +goose Up
-- Migration: 003_projects.sql
-- Purpose: Projects tracking and project milestones.

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active', -- 'active' | 'completed' | 'archived' | 'on_hold'
    color TEXT NOT NULL DEFAULT '#3b82f6',
    progress_percentage INTEGER NOT NULL DEFAULT 0,
    target_date DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS milestones (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    due_date DATETIME,
    is_completed BOOLEAN NOT NULL DEFAULT 0,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones(project_id);

-- +goose Down
DROP TABLE IF EXISTS milestones;
DROP TABLE IF EXISTS projects;
