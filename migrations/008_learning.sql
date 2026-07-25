-- +goose Up
-- Migration: 008_learning.sql
-- Purpose: Learning tracks and study session logs (v2 support).

CREATE TABLE IF NOT EXISTS learning_tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    category TEXT NOT NULL, -- 'dsa' | 'course' | 'book' | 'paper' | 'language' | 'framework'
    description TEXT NOT NULL DEFAULT '',
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS learning_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    track_id INTEGER NOT NULL REFERENCES learning_tracks(id) ON DELETE CASCADE,
    duration_minutes INTEGER NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    session_date DATE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_learning_tracks_user ON learning_tracks(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_track ON learning_sessions(track_id);

-- +goose Down
DROP TABLE IF EXISTS learning_sessions;
DROP TABLE IF EXISTS learning_tracks;
