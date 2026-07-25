package integration

import (
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
)

func TestDB_MigrateUpAndTables(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "integration_test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	// Apply all 9 embedded migrations
	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("MigrateUp failed: %v", err)
	}

	// Verify required tables exist
	expectedTables := []string{
		"users",
		"sessions",
		"settings",
		"projects",
		"milestones",
		"tasks",
		"task_labels",
		"task_dependencies",
		"notes",
		"note_tags",
		"note_links",
		"journal_entries",
		"journal_items",
		"transactions",
		"budgets",
		"learning_tracks",
		"learning_sessions",
		"search_index",
	}

	for _, table := range expectedTables {
		var name string
		err := database.Raw.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?;", table).Scan(&name)
		if err != nil {
			t.Errorf("expected table '%s' to exist after migration, but query failed: %v", table, err)
		}
	}
}
