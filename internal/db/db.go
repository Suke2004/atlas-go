// Package db wraps SQLite database connection initialization, WAL mode configuration,
// and transaction management for Atlas.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps *sql.DB and embeds the sqlc-generated Queries struct.
type DB struct {
	*Queries
	Raw *sql.DB
}

// Open opens a connection to the SQLite database at dbPath, creates parent directories if needed,
// and enforces high-performance, WAL-mode pragmas.
func Open(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("db: failed to create directory %s: %w", dir, err)
	}

	raw, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("db: failed to open sqlite database: %w", err)
	}

	raw.SetMaxOpenConns(1) // SQLite write lock safety
	raw.SetMaxIdleConns(1)
	raw.SetConnMaxLifetime(1 * time.Hour)

	if err := raw.Ping(); err != nil {
		return nil, fmt.Errorf("db: failed to ping sqlite database: %w", err)
	}

	// Apply WAL mode and essential pragmas explicitly
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA synchronous=NORMAL;",
	}
	for _, pragma := range pragmas {
		if _, err := raw.Exec(pragma); err != nil {
			return nil, fmt.Errorf("db: failed to execute pragma '%s': %w", pragma, err)
		}
	}

	return &DB{
		Queries: New(raw),
		Raw:     raw,
	}, nil
}

// Close closes the underlying raw database connection.
func (d *DB) Close() error {
	if d.Raw != nil {
		return d.Raw.Close()
	}
	return nil
}

// WithTx runs fn within an explicit SQL transaction.
// If fn returns an error, the transaction is rolled back. Otherwise it is committed.
func (d *DB) WithTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := d.Raw.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db: begin tx: %w", err)
	}

	q := New(tx)
	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("db: tx err: %v, rollback err: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db: commit tx: %w", err)
	}

	return nil
}

// WithRawTx runs fn within a raw SQL transaction when raw SQL statements are required.
func (d *DB) WithRawTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := d.Raw.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db: begin raw tx: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("db: raw tx err: %v, rollback err: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("db: commit raw tx: %w", err)
	}

	return nil
}
