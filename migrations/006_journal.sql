-- +goose Up
-- Migration: 006_journal.sql
-- Purpose: Daily journal entries, mood/energy/sleep metrics, and daily item bullet points.

CREATE TABLE IF NOT EXISTS journal_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entry_date DATE NOT NULL,
    mood_rating INTEGER CHECK(mood_rating BETWEEN 1 AND 5),
    energy_rating INTEGER CHECK(energy_rating BETWEEN 1 AND 5),
    sleep_hours REAL,
    summary TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, entry_date)
);

CREATE TABLE IF NOT EXISTS journal_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_id INTEGER NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    category TEXT NOT NULL, -- 'win' | 'problem' | 'idea' | 'tomorrow'
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_journal_entries_user_date ON journal_entries(user_id, entry_date);
CREATE INDEX IF NOT EXISTS idx_journal_items_entry ON journal_items(entry_id);

-- +goose Down
DROP TABLE IF EXISTS journal_items;
DROP TABLE IF EXISTS journal_entries;
