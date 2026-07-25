-- +goose Up
-- Migration: 005_notes.sql
-- Purpose: Knowledge base notes, note tags, and bidirectional links between notes.

CREATE TABLE IF NOT EXISTS notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    is_pinned BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS note_tags (
    note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    PRIMARY KEY (note_id, tag)
);

CREATE TABLE IF NOT EXISTS note_links (
    source_note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    target_note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    PRIMARY KEY (source_note_id, target_note_id)
);

CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_project_id ON notes(project_id);
CREATE INDEX IF NOT EXISTS idx_note_tags_tag ON note_tags(tag);

-- +goose Down
DROP TABLE IF EXISTS note_links;
DROP TABLE IF EXISTS note_tags;
DROP TABLE IF EXISTS notes;
