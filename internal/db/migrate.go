// Package db provides programmatic database migration capabilities using Goose.
package db

import (
	"database/sql"
	"fmt"

	"github.com/Suke2004/atlas-go/migrations"
	"github.com/pressly/goose/v3"
)

// MigrateUp applies all pending Goose SQL migrations to the SQLite database.
func MigrateUp(database *sql.DB) error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("migrate: failed to set dialect: %w", err)
	}

	if err := goose.Up(database, "."); err != nil {
		return fmt.Errorf("migrate: failed to apply migrations: %w", err)
	}

	return nil
}

// MigrateStatus returns the current goose migration status.
func MigrateStatus(database *sql.DB) error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("migrate: failed to set dialect: %w", err)
	}

	return goose.Status(database, ".")
}

