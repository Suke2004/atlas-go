-- +goose Up
-- Documents table: stores file metadata. Binary content lives in /data/uploads/.

CREATE TABLE IF NOT EXISTS documents (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename      TEXT    NOT NULL,           -- stored filename (uuid.ext)
    original_name TEXT    NOT NULL,           -- original upload filename
    mime_type     TEXT    NOT NULL DEFAULT '',
    file_size     INTEGER NOT NULL DEFAULT 0, -- bytes
    storage_path  TEXT    NOT NULL,           -- relative to /data/uploads/
    title         TEXT    NOT NULL DEFAULT '',
    summary       TEXT,                       -- AI-generated or manual
    content_text  TEXT,                       -- extracted text for FTS5
    tags          TEXT    NOT NULL DEFAULT '[]', -- JSON array of strings
    created_at    DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_documents_user_id   ON documents(user_id);
CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents(created_at DESC);

-- FTS5 entry for documents
INSERT OR IGNORE INTO search_index(entity_type, entity_id, user_id, title, content, tags)
SELECT 'document', id, user_id, original_name, COALESCE(content_text, ''), tags
FROM documents
WHERE 1=0; -- empty seed; triggers below keep it live

-- Trigger: insert into FTS on new document
-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS fts_documents_insert
AFTER INSERT ON documents
BEGIN
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('document', NEW.id, NEW.user_id, NEW.original_name, COALESCE(NEW.content_text, ''), NEW.tags);
END;
-- +goose StatementEnd

-- Trigger: update FTS on document edit
-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS fts_documents_update
AFTER UPDATE ON documents
BEGIN
    DELETE FROM search_index WHERE entity_type = 'document' AND entity_id = OLD.id;
    INSERT INTO search_index(entity_type, entity_id, user_id, title, content, tags)
    VALUES ('document', NEW.id, NEW.user_id, NEW.original_name, COALESCE(NEW.content_text, ''), NEW.tags);
END;
-- +goose StatementEnd

-- Trigger: remove from FTS on document delete
-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS fts_documents_delete
AFTER DELETE ON documents
BEGIN
    DELETE FROM search_index WHERE entity_type = 'document' AND entity_id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS fts_documents_delete;
DROP TRIGGER IF EXISTS fts_documents_update;
DROP TRIGGER IF EXISTS fts_documents_insert;
DROP INDEX  IF EXISTS idx_documents_created_at;
DROP INDEX  IF EXISTS idx_documents_user_id;
DROP TABLE  IF EXISTS documents;
