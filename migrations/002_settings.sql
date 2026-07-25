-- +goose Up
-- Migration: 002_settings.sql
-- Purpose: Key-value configuration store per user (theme, widget layout, default preferences).

CREATE TABLE IF NOT EXISTS settings (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, key)
);

-- +goose Down
DROP TABLE IF EXISTS settings;
