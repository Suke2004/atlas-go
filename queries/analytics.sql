-- name: GetDailyTaskCounts :many
SELECT DATE(updated_at) AS date_str, COUNT(*) as count
FROM tasks
WHERE user_id = ? AND status = 'completed' AND updated_at >= ?
GROUP BY DATE(updated_at);

-- name: GetDailyNoteCounts :many
SELECT DATE(created_at) AS date_str, COUNT(*) as count
FROM notes
WHERE user_id = ? AND created_at >= ?
GROUP BY DATE(created_at);

-- name: GetDailyJournalCounts :many
SELECT entry_date AS date_str, COUNT(*) as count
FROM journal_entries
WHERE user_id = ? AND entry_date >= ?
GROUP BY entry_date;

-- name: GetDailyLearningCounts :many
SELECT session_date AS date_str, COUNT(*) as count
FROM learning_sessions ls
JOIN learning_tracks lt ON ls.track_id = lt.id
WHERE lt.user_id = ? AND ls.session_date >= ?
GROUP BY ls.session_date;

-- name: GetDailyTransactionCounts :many
SELECT transaction_date AS date_str, COUNT(*) as count
FROM transactions
WHERE user_id = ? AND transaction_date >= ?
GROUP BY transaction_date;

-- name: GetDailyDocumentCounts :many
SELECT DATE(created_at) AS date_str, COUNT(*) as count
FROM documents
WHERE user_id = ? AND created_at >= ?
GROUP BY DATE(created_at);

-- name: GetCategoryExpensesCurrentMonth :many
SELECT category, SUM(amount) as total_amount
FROM transactions
WHERE user_id = ? AND type = 'expense' AND strftime('%Y-%m', transaction_date) = strftime('%Y-%m', 'now')
GROUP BY category;

-- name: GetMoodEnergyTrends30Days :many
SELECT entry_date, mood_score, energy_rating, sleep_hours
FROM journal_entries
WHERE user_id = ? AND entry_date >= date('now', '-30 days')
ORDER BY entry_date ASC;
