package db

import (
	"context"
)

const createNote = `-- name: CreateNote :one
INSERT INTO notes (user_id, project_id, title, content, is_pinned)
VALUES (?, ?, ?, ?, ?)
RETURNING id, user_id, project_id, title, content, is_pinned, created_at, updated_at
`

type CreateNoteParams struct {
	UserID    int64       `json:"user_id"`
	ProjectID interface{} `json:"project_id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	IsPinned  bool        `json:"is_pinned"`
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.UserID,
		arg.ProjectID,
		arg.Title,
		arg.Content,
		arg.IsPinned,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProjectID,
		&i.Title,
		&i.Content,
		&i.IsPinned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listNotes = `-- name: ListNotes :many
SELECT id, user_id, project_id, title, content, is_pinned, created_at, updated_at FROM notes
WHERE user_id = ?
ORDER BY is_pinned DESC, updated_at DESC
`

func (q *Queries) ListNotes(ctx context.Context, userID int64) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, listNotes, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProjectID,
			&i.Title,
			&i.Content,
			&i.IsPinned,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const addNoteTag = `-- name: AddNoteTag :exec
INSERT INTO note_tags (note_id, tag)
VALUES (?, ?)
ON CONFLICT (note_id, tag) DO NOTHING
`

type AddNoteTagParams struct {
	NoteID int64  `json:"note_id"`
	Tag    string `json:"tag"`
}

func (q *Queries) AddNoteTag(ctx context.Context, arg AddNoteTagParams) error {
	_, err := q.db.ExecContext(ctx, addNoteTag, arg.NoteID, arg.Tag)
	return err
}
