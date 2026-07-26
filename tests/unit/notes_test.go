package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestNotes_ServiceWikiLinksAndTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "notes_test.db")

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
		Username:    "noteuser",
		DisplayName: "Note User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	notesRepo := notes.NewRepository(database)
	notesSvc := notes.NewService(notesRepo, nil)

	// Create Target Note
	targetNote, err := notesSvc.CreateNote(ctx, user.ID, notes.NoteInput{
		Title:   "Architecture Overview",
		Content: "System architecture details.",
		Tags:    []string{"architecture", "docs"},
	})
	if err != nil {
		t.Fatalf("failed to create target note: %v", err)
	}

	// Create Source Note referencing [[Architecture Overview]]
	_, err = notesSvc.CreateNote(ctx, user.ID, notes.NoteInput{
		Title:   "Sprint 1 Meeting",
		Content: "Discussed system design based on [[Architecture Overview]].",
		Tags:    []string{"meeting"},
	})
	if err != nil {
		t.Fatalf("failed to create source note: %v", err)
	}

	// Verify Backlink on Target Note
	targetDetails, err := notesSvc.GetNote(ctx, user.ID, targetNote.ID)
	if err != nil {
		t.Fatalf("failed to get target note details: %v", err)
	}

	if len(targetDetails.Backlinks) != 1 {
		t.Fatalf("expected 1 backlink on target note, got: %d", len(targetDetails.Backlinks))
	}

	if targetDetails.Backlinks[0].Title != "Sprint 1 Meeting" {
		t.Errorf("expected backlink title 'Sprint 1 Meeting', got: %s", targetDetails.Backlinks[0].Title)
	}

	summary, err := notesSvc.GetNotesSummary(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get notes summary: %v", err)
	}

	if summary.TotalNotes != 2 {
		t.Errorf("expected 2 total notes, got: %d", summary.TotalNotes)
	}
}
