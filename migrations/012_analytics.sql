-- +goose Up
-- Migration: 012_analytics.sql
-- Purpose: Indexes to optimize activity heatmap aggregation and analytics reporting.

CREATE INDEX IF NOT EXISTS idx_tasks_completed_at ON tasks(updated_at) WHERE status = 'completed';
CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);
CREATE INDEX IF NOT EXISTS idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_date ON learning_sessions(session_date);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(transaction_date);

-- +goose Down
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_learning_sessions_date;
DROP INDEX IF EXISTS idx_journal_entries_date;
DROP INDEX IF EXISTS idx_notes_created_at;
DROP INDEX IF EXISTS idx_tasks_completed_at;
