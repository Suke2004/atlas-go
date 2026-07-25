-- +goose Up
-- Migration: 007_finance.sql
-- Purpose: Financial transaction tracking and monthly budget allocation (v2 support).

CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount REAL NOT NULL,
    type TEXT NOT NULL, -- 'income' | 'expense'
    category TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    transaction_date DATE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS budgets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category TEXT NOT NULL,
    monthly_limit REAL NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, category)
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, transaction_date);

-- +goose Down
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS transactions;
