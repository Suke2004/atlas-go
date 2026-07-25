-- name: CreateUser :one
INSERT INTO users (username, display_name, email, password_hash, timezone)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateSession :one
INSERT INTO sessions (id, user_id, expires_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = ? AND expires_at > CURRENT_TIMESTAMP LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP;
