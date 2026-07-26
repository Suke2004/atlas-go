-- name: CreateDocument :one
INSERT INTO documents (user_id, filename, original_name, mime_type, file_size, storage_path, title, tags)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetDocument :one
SELECT * FROM documents
WHERE id = ? AND user_id = ?
LIMIT 1;

-- name: ListDocuments :many
SELECT * FROM documents
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: UpdateDocumentMeta :one
UPDATE documents
SET title      = ?,
    summary    = ?,
    tags       = ?,
    updated_at = datetime('now')
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateDocumentContent :exec
UPDATE documents
SET content_text = ?,
    updated_at   = datetime('now')
WHERE id = ? AND user_id = ?;

-- name: UpdateDocumentSummary :exec
UPDATE documents
SET summary    = ?,
    updated_at = datetime('now')
WHERE id = ? AND user_id = ?;

-- name: DeleteDocument :exec
DELETE FROM documents WHERE id = ? AND user_id = ?;

-- name: CountDocuments :one
SELECT COUNT(*) FROM documents WHERE user_id = ?;

-- name: GetDocumentsByTag :many
SELECT * FROM documents
WHERE user_id = ? AND tags LIKE ?
ORDER BY created_at DESC;
