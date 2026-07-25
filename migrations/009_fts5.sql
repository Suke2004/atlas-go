-- +goose Up
-- Migration: 009_fts5.sql
-- Purpose: Global full-text search index (FTS5) and automatic triggers to sync
-- tasks, notes, projects, and journal entries.

CREATE VIRTUAL TABLE IF NOT EXISTS search_index USING fts5(
    entity_type,          -- 'task' | 'note' | 'project' | 'journal'
    entity_id UNINDEXED,  -- ID in source table
    user_id UNINDEXED,    -- Owner ID
    title,
    content,
    tags
);

-- ── 1. Tasks Triggers ───────────────────────────────────────────────────────

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_insert
AFTER INSERT ON tasks
BEGIN
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('task', NEW.id, NEW.user_id, NEW.title, NEW.description, NEW.priority);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_update
AFTER UPDATE ON tasks
BEGIN
    DELETE FROM search_index WHERE entity_type = 'task' AND entity_id = OLD.id;
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('task', NEW.id, NEW.user_id, NEW.title, NEW.description, NEW.priority);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_tasks_fts_delete
AFTER DELETE ON tasks
BEGIN
    DELETE FROM search_index WHERE entity_type = 'task' AND entity_id = OLD.id;
END;
-- +goose StatementEnd

-- ── 2. Notes Triggers ───────────────────────────────────────────────────────

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_notes_fts_insert
AFTER INSERT ON notes
BEGIN
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('note', NEW.id, NEW.user_id, NEW.title, NEW.content, '');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_notes_fts_update
AFTER UPDATE ON notes
BEGIN
    DELETE FROM search_index WHERE entity_type = 'note' AND entity_id = OLD.id;
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('note', NEW.id, NEW.user_id, NEW.title, NEW.content, '');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_notes_fts_delete
AFTER DELETE ON notes
BEGIN
    DELETE FROM search_index WHERE entity_type = 'note' AND entity_id = OLD.id;
END;
-- +goose StatementEnd

-- ── 3. Projects Triggers ────────────────────────────────────────────────────

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_projects_fts_insert
AFTER INSERT ON projects
BEGIN
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('project', NEW.id, NEW.user_id, NEW.name, NEW.description, NEW.status);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_projects_fts_update
AFTER UPDATE ON projects
BEGIN
    DELETE FROM search_index WHERE entity_type = 'project' AND entity_id = OLD.id;
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('project', NEW.id, NEW.user_id, NEW.name, NEW.description, NEW.status);
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_projects_fts_delete
AFTER DELETE ON projects
BEGIN
    DELETE FROM search_index WHERE entity_type = 'project' AND entity_id = OLD.id;
END;
-- +goose StatementEnd

-- ── 4. Journal Triggers ────────────────────────────────────────────────────

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_journal_fts_insert
AFTER INSERT ON journal_entries
BEGIN
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('journal', NEW.id, NEW.user_id, 'Journal Entry ' || NEW.entry_date, NEW.summary, '');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_journal_fts_update
AFTER UPDATE ON journal_entries
BEGIN
    DELETE FROM search_index WHERE entity_type = 'journal' AND entity_id = OLD.id;
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('journal', NEW.id, NEW.user_id, 'Journal Entry ' || NEW.entry_date, NEW.summary, '');
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS trg_journal_fts_delete
AFTER DELETE ON journal_entries
BEGIN
    DELETE FROM search_index WHERE entity_type = 'journal' AND entity_id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_journal_fts_delete;
DROP TRIGGER IF EXISTS trg_journal_fts_update;
DROP TRIGGER IF EXISTS trg_journal_fts_insert;
DROP TRIGGER IF EXISTS trg_projects_fts_delete;
DROP TRIGGER IF EXISTS trg_projects_fts_update;
DROP TRIGGER IF EXISTS trg_projects_fts_insert;
DROP TRIGGER IF EXISTS trg_notes_fts_delete;
DROP TRIGGER IF EXISTS trg_notes_fts_update;
DROP TRIGGER IF EXISTS trg_notes_fts_insert;
DROP TRIGGER IF EXISTS trg_tasks_fts_delete;
DROP TRIGGER IF EXISTS trg_tasks_fts_update;
DROP TRIGGER IF EXISTS trg_tasks_fts_insert;
DROP TABLE IF EXISTS search_index;
