// Code generated manually matching sqlc conventions for analytics queries.
package db

import (
	"context"
)

const getDailyTaskCounts = `-- name: GetDailyTaskCounts :many
SELECT DATE(updated_at) AS date_str, COUNT(*) as count
FROM tasks
WHERE user_id = ? AND status = 'completed' AND updated_at >= ?
GROUP BY DATE(updated_at)
`

type DailyCountRow struct {
	DateStr string `json:"date_str"`
	Count   int64  `json:"count"`
}

func (q *Queries) GetDailyTaskCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyTaskCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getDailyNoteCounts = `-- name: GetDailyNoteCounts :many
SELECT DATE(created_at) AS date_str, COUNT(*) as count
FROM notes
WHERE user_id = ? AND created_at >= ?
GROUP BY DATE(created_at)
`

func (q *Queries) GetDailyNoteCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyNoteCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getDailyJournalCounts = `-- name: GetDailyJournalCounts :many
SELECT entry_date AS date_str, COUNT(*) as count
FROM journal_entries
WHERE user_id = ? AND entry_date >= ?
GROUP BY entry_date
`

func (q *Queries) GetDailyJournalCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyJournalCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getDailyLearningCounts = `-- name: GetDailyLearningCounts :many
SELECT session_date AS date_str, COUNT(*) as count
FROM learning_sessions ls
JOIN learning_tracks lt ON ls.track_id = lt.id
WHERE lt.user_id = ? AND ls.session_date >= ?
GROUP BY ls.session_date
`

func (q *Queries) GetDailyLearningCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyLearningCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getDailyTransactionCounts = `-- name: GetDailyTransactionCounts :many
SELECT transaction_date AS date_str, COUNT(*) as count
FROM transactions
WHERE user_id = ? AND transaction_date >= ?
GROUP BY transaction_date
`

func (q *Queries) GetDailyTransactionCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyTransactionCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getDailyDocumentCounts = `-- name: GetDailyDocumentCounts :many
SELECT DATE(created_at) AS date_str, COUNT(*) as count
FROM documents
WHERE user_id = ? AND created_at >= ?
GROUP BY DATE(created_at)
`

func (q *Queries) GetDailyDocumentCounts(ctx context.Context, userID int64, since string) ([]DailyCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getDailyDocumentCounts, userID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyCountRow
	for rows.Next() {
		var i DailyCountRow
		if err := rows.Scan(&i.DateStr, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getCategoryExpensesCurrentMonth = `-- name: GetCategoryExpensesCurrentMonth :many
SELECT category, SUM(amount) as total_amount
FROM transactions
WHERE user_id = ? AND type = 'expense' AND strftime('%Y-%m', transaction_date) = strftime('%Y-%m', 'now')
GROUP BY category
`

type CategoryExpenseRow struct {
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
}

func (q *Queries) GetCategoryExpensesCurrentMonth(ctx context.Context, userID int64) ([]CategoryExpenseRow, error) {
	rows, err := q.db.QueryContext(ctx, getCategoryExpensesCurrentMonth, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoryExpenseRow
	for rows.Next() {
		var i CategoryExpenseRow
		if err := rows.Scan(&i.Category, &i.TotalAmount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getMoodEnergyTrends30Days = `-- name: GetMoodEnergyTrends30Days :many
SELECT entry_date, mood_score, energy_rating, sleep_hours
FROM journal_entries
WHERE user_id = ? AND entry_date >= date('now', '-30 days')
ORDER BY entry_date ASC
`

type MoodEnergyRow struct {
	EntryDate    string  `json:"entry_date"`
	MoodScore    int64   `json:"mood_score"`
	EnergyRating int64   `json:"energy_rating"`
	SleepHours   float64 `json:"sleep_hours"`
}

func (q *Queries) GetMoodEnergyTrends30Days(ctx context.Context, userID int64) ([]MoodEnergyRow, error) {
	rows, err := q.db.QueryContext(ctx, getMoodEnergyTrends30Days, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MoodEnergyRow
	for rows.Next() {
		var i MoodEnergyRow
		if err := rows.Scan(&i.EntryDate, &i.MoodScore, &i.EnergyRating, &i.SleepHours); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
