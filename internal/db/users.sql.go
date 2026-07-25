package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, display_name, email, password_hash, timezone)
VALUES (?, ?, ?, ?, ?)
RETURNING id, username, display_name, email, password_hash, timezone, created_at, updated_at
`

type CreateUserParams struct {
	Username     string      `json:"username"`
	DisplayName  string      `json:"display_name"`
	Email        interface{} `json:"email"`
	PasswordHash string      `json:"password_hash"`
	Timezone     string      `json:"timezone"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.DisplayName,
		arg.Email,
		arg.PasswordHash,
		arg.Timezone,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.Timezone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, display_name, email, password_hash, timezone, created_at, updated_at FROM users
WHERE id = ? LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.Timezone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, display_name, email, password_hash, timezone, created_at, updated_at FROM users
WHERE username = ? LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.DisplayName,
		&i.Email,
		&i.PasswordHash,
		&i.Timezone,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES (?, ?, ?)
RETURNING id, user_id, expires_at, created_at
`

type CreateSessionParams struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession, arg.ID, arg.UserID, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionByID = `-- name: GetSessionByID :one
SELECT id, user_id, expires_at, created_at FROM sessions
WHERE id = ? AND expires_at > CURRENT_TIMESTAMP LIMIT 1
`

func (q *Queries) GetSessionByID(ctx context.Context, id string) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByID, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?
`

func (q *Queries) DeleteSession(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteSession, id)
	return err
}

const deleteExpiredSessions = `-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP
`

func (q *Queries) DeleteExpiredSessions(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteExpiredSessions)
	return err
}
