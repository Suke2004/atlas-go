-- name: CreateTransaction :one
INSERT INTO transactions (user_id, amount, type, category, description, transaction_date)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE user_id = ?
ORDER BY transaction_date DESC, created_at DESC;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = ? AND user_id = ?;

-- name: GetFinanceSummary :one
SELECT 
    COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0.0) AS total_income,
    COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0.0) AS total_expenses
FROM transactions
WHERE user_id = ?;
