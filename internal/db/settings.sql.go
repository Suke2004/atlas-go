package db

import (
	"context"
)

const getSetting = `-- name: GetSetting :one
SELECT value FROM settings
WHERE user_id = ? AND key = ? LIMIT 1
`

type GetSettingParams struct {
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
}

func (q *Queries) GetSetting(ctx context.Context, arg GetSettingParams) (string, error) {
	row := q.db.QueryRowContext(ctx, getSetting, arg.UserID, arg.Key)
	var value string
	err := row.Scan(&value)
	return value, err
}

const getAllSettings = `-- name: GetAllSettings :many
SELECT key, value FROM settings
WHERE user_id = ?
`

type SettingRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (q *Queries) GetAllSettings(ctx context.Context, userID int64) ([]SettingRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllSettings, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SettingRow
	for rows.Next() {
		var i SettingRow
		if err := rows.Scan(&i.Key, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setSetting = `-- name: SetSetting :exec
INSERT INTO settings (user_id, key, value, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
ON CONFLICT (user_id, key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = CURRENT_TIMESTAMP
`

type SetSettingParams struct {
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (q *Queries) SetSetting(ctx context.Context, arg SetSettingParams) error {
	_, err := q.db.ExecContext(ctx, setSetting, arg.UserID, arg.Key, arg.Value)
	return err
}

const deleteSetting = `-- name: DeleteSetting :exec
DELETE FROM settings
WHERE user_id = ? AND key = ?
`

type DeleteSettingParams struct {
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
}

func (q *Queries) DeleteSetting(ctx context.Context, arg DeleteSettingParams) error {
	_, err := q.db.ExecContext(ctx, deleteSetting, arg.UserID, arg.Key)
	return err
}
