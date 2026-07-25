-- name: CreateNote :one
INSERT INTO notes (user_id, project_id, title, content, is_pinned)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetNoteByID :one
SELECT * FROM notes
WHERE id = ? AND user_id = ? LIMIT 1;

-- name: ListNotes :many
SELECT * FROM notes
WHERE user_id = ?
ORDER BY is_pinned DESC, updated_at DESC;

-- name: UpdateNote :one
UPDATE notes
SET project_id = ?, title = ?, content = ?, is_pinned = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = ? AND user_id = ?;

-- name: AddNoteTag :exec
INSERT INTO note_tags (note_id, tag)
VALUES (?, ?)
ON CONFLICT (note_id, tag) DO NOTHING;

-- name: ListNoteTags :many
SELECT tag FROM note_tags
WHERE note_id = ?;

-- name: AddNoteLink :exec
INSERT INTO note_links (source_note_id, target_note_id)
VALUES (?, ?)
ON CONFLICT (source_note_id, target_note_id) DO NOTHING;

-- name: GetNoteBacklinks :many
SELECT n.* FROM notes n
JOIN note_links nl ON n.id = nl.source_note_id
WHERE nl.target_note_id = ? AND n.user_id = ?;
