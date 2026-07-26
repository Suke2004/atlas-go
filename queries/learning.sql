-- name: CreateLearningTrack :one
INSERT INTO learning_tracks (user_id, title, category, description, current_streak, longest_streak)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListLearningTracks :many
SELECT * FROM learning_tracks
WHERE user_id = ?
ORDER BY updated_at DESC;

-- name: DeleteLearningTrack :exec
DELETE FROM learning_tracks
WHERE id = ? AND user_id = ?;

-- name: AddLearningSession :one
INSERT INTO learning_sessions (track_id, duration_minutes, summary, session_date)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: ListLearningSessions :many
SELECT * FROM learning_sessions
WHERE track_id = ?
ORDER BY session_date DESC;
