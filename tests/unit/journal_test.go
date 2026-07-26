package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/journal"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestJournal_ServiceMindSyncAndItems(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "journal_test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	ctx := context.Background()
	setupSvc := setup.NewService(database)
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "journaluser",
		DisplayName: "Journal User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	journalRepo := journal.NewRepository(database)
	journalSvc := journal.NewService(journalRepo, nil, nil)

	// Save Daily Entry
	entry, err := journalSvc.SaveDailyJournal(ctx, user.ID, journal.JournalInput{
		EntryDate:    "2026-07-26",
		MoodRating:   5,
		EnergyRating: 4,
		SleepHours:  8.0,
		Summary:      "Great productivity day with deep focus.",
	})
	if err != nil {
		t.Fatalf("failed to save daily journal: %v", err)
	}

	// Add 4-Quadrant Items
	winItem, err := journalSvc.AddJournalItem(ctx, user.ID, entry.ID, "win", "Shipped Phase 8 and Phase 9")
	if err != nil {
		t.Fatalf("failed to add win item: %v", err)
	}

	if winItem.Content != "Shipped Phase 8 and Phase 9" {
		t.Errorf("expected item content 'Shipped Phase 8 and Phase 9', got: %s", winItem.Content)
	}

	// Verify Daily Details
	details, err := journalSvc.GetDailyJournal(ctx, user.ID, "2026-07-26")
	if err != nil {
		t.Fatalf("failed to get daily journal details: %v", err)
	}

	if len(details.Wins) != 1 {
		t.Fatalf("expected 1 win item, got: %d", len(details.Wins))
	}
}
