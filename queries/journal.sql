-- name: UpsertJournalEntry :one
INSERT INTO journal_entries (user_id, entry_date, mood_rating, energy_rating, sleep_hours, summary)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (user_id, entry_date) DO UPDATE SET
    mood_rating = EXCLUDED.mood_rating,
    energy_rating = EXCLUDED.energy_rating,
    sleep_hours = EXCLUDED.sleep_hours,
    summary = EXCLUDED.summary,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetJournalEntryByDate :one
SELECT * FROM journal_entries
WHERE user_id = ? AND entry_date = ? LIMIT 1;

-- name: ListJournalEntries :many
SELECT * FROM journal_entries
WHERE user_id = ?
ORDER BY entry_date DESC
LIMIT 30;

-- name: AddJournalItem :one
INSERT INTO journal_items (entry_id, category, content)
VALUES (?, ?, ?)
RETURNING *;

-- name: ListJournalItems :many
SELECT * FROM journal_items
WHERE entry_id = ?
ORDER BY created_at ASC;

-- name: DeleteJournalItem :exec
DELETE FROM journal_items
WHERE id = ?;
