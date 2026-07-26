-- +goose Up
-- Migration: 010_github_projects.sql
-- Purpose: Add GitHub integration and tech stack columns to projects table.

ALTER TABLE projects ADD COLUMN github_url TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN github_stars INTEGER NOT NULL DEFAULT 0;
ALTER TABLE projects ADD COLUMN github_forks INTEGER NOT NULL DEFAULT 0;
ALTER TABLE projects ADD COLUMN github_open_issues INTEGER NOT NULL DEFAULT 0;
ALTER TABLE projects ADD COLUMN github_language TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD COLUMN github_last_pushed_at DATETIME;
ALTER TABLE projects ADD COLUMN tech_stack TEXT NOT NULL DEFAULT '';

-- +goose Down
-- SQLite does not support ALTER TABLE DROP COLUMN directly in older syntax, statement ignored on down.
