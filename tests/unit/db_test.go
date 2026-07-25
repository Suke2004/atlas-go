package unit

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
)

func TestDB_OpenPragmas(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_atlas.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer database.Close()

	// Verify WAL mode pragma
	var journalMode string
	err = database.Raw.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
	if err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}

	if journalMode != "wal" {
		t.Errorf("expected journal_mode 'wal', got '%s'", journalMode)
	}

	// Verify Foreign Keys pragma
	var foreignKeys int
	err = database.Raw.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("failed to query foreign_keys: %v", err)
	}

	if foreignKeys != 1 {
		t.Errorf("expected foreign_keys 1, got %d", foreignKeys)
	}
}

func TestDB_WithTx_RollbackOnFailure(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_tx.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	defer database.Close()

	// Create dummy table for transaction test
	_, err = database.Raw.Exec("CREATE TABLE test_items (id INT PRIMARY KEY, name TEXT);")
	if err != nil {
		t.Fatalf("failed to create dummy table: %v", err)
	}

	// Run transaction that executes an insert on tx and returns error
	ctx := context.Background()
	_ = database.WithRawTx(ctx, func(tx *sql.Tx) error {
		_, execErr := tx.Exec("INSERT INTO test_items (id, name) VALUES (1, 'item1');")
		if execErr != nil {
			return execErr
		}
		return context.DeadlineExceeded // simulate error
	})

	// Verify item was not committed
	var count int
	_ = database.Raw.QueryRow("SELECT COUNT(*) FROM test_items;").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 items after rolled back transaction, got %d", count)
	}
}
