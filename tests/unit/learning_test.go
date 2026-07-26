package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/learning"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestLearning_ServiceAndSessions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "learning_test.db")

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
		Username:    "learnuser",
		DisplayName: "Learn User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	learnRepo := learning.NewRepository(database)
	learnSvc := learning.NewService(learnRepo)

	// Create track
	track, err := learnSvc.CreateTrack(ctx, user.ID, learning.TrackInput{
		Title:       "Go Concurrency Patterns",
		Category:    "language",
		Description: "Channels and Goroutines",
	})
	if err != nil {
		t.Fatalf("failed to create track: %v", err)
	}

	// Add 60-min session
	_, err = learnSvc.AddSession(ctx, user.ID, track.ID, 60, "Completed worker pool implementation")
	if err != nil {
		t.Fatalf("failed to add session: %v", err)
	}

	summary, err := learnSvc.GetLearningSummary(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}

	if summary.TotalTracks != 1 {
		t.Errorf("expected 1 track, got: %d", summary.TotalTracks)
	}
	if summary.TotalStudyHours != 1 {
		t.Errorf("expected 1 study hour, got: %d", summary.TotalStudyHours)
	}
}
