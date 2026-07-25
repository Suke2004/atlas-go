-- name: GetSetting :one
SELECT value FROM settings
WHERE user_id = ? AND key = ? LIMIT 1;

-- name: GetAllSettings :many
SELECT key, value FROM settings
WHERE user_id = ?;

-- name: SetSetting :exec
INSERT INTO settings (user_id, key, value, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT (user_id, key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = CURRENT_TIMESTAMP;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE user_id = ? AND key = ?;
