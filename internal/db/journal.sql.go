package db

import (
	"context"
)

const upsertJournalEntry = `-- name: UpsertJournalEntry :one
INSERT INTO journal_entries (user_id, entry_date, mood_rating, energy_rating, sleep_hours, summary)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (user_id, entry_date) DO UPDATE SET
    mood_rating = EXCLUDED.mood_rating,
    energy_rating = EXCLUDED.energy_rating,
    sleep_hours = EXCLUDED.sleep_hours,
    summary = EXCLUDED.summary,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, user_id, entry_date, mood_rating, energy_rating, sleep_hours, summary, created_at, updated_at
`

type UpsertJournalEntryParams struct {
	UserID       int64       `json:"user_id"`
	EntryDate    string      `json:"entry_date"`
	MoodRating   interface{} `json:"mood_rating"`
	EnergyRating interface{} `json:"energy_rating"`
	SleepHours  interface{} `json:"sleep_hours"`
	Summary      string      `json:"summary"`
}

func (q *Queries) UpsertJournalEntry(ctx context.Context, arg UpsertJournalEntryParams) (JournalEntry, error) {
	row := q.db.QueryRowContext(ctx, upsertJournalEntry,
		arg.UserID,
		arg.EntryDate,
		arg.MoodRating,
		arg.EnergyRating,
		arg.SleepHours,
		arg.Summary,
	)
	var i JournalEntry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EntryDate,
		&i.MoodRating,
		&i.EnergyRating,
		&i.SleepHours,
		&i.Summary,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addJournalItem = `-- name: AddJournalItem :one
INSERT INTO journal_items (entry_id, category, content)
VALUES (?, ?, ?)
RETURNING id, entry_id, category, content, created_at
`

type AddJournalItemParams struct {
	EntryID  int64  `json:"entry_id"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func (q *Queries) AddJournalItem(ctx context.Context, arg AddJournalItemParams) (JournalItem, error) {
	row := q.db.QueryRowContext(ctx, addJournalItem, arg.EntryID, arg.Category, arg.Content)
	var i JournalItem
	err := row.Scan(
		&i.ID,
		&i.EntryID,
		&i.Category,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}
