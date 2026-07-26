package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/search"
	"github.com/Suke2004/atlas-go/internal/setup"
	"github.com/Suke2004/atlas-go/internal/tasks"
)

func TestSearch_GlobalSearch(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "search_test.db")

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
		Username:    "searchuser",
		DisplayName: "Search User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	projRepo := projects.NewRepository(database)
	projSvc := projects.NewService(projRepo, nil)

	tasksRepo := tasks.NewRepository(database)
	tasksSvc := tasks.NewService(tasksRepo, projRepo, projSvc)

	notesRepo := notes.NewRepository(database)
	notesSvc := notes.NewService(notesRepo, projRepo)

	// Create test project and note
	_, err = projSvc.CreateProject(ctx, user.ID, projects.ProjectInput{
		Name:        "Atlas Operating System",
		Description: "Self-hosted workspace platform.",
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	n, err := notesSvc.CreateNote(ctx, user.ID, notes.NoteInput{
		Title:   "Atlas Architecture Guide",
		Content: "Layered architecture design.",
	})
	if err != nil {
		t.Fatalf("failed to create note: %v", err)
	}
	searchSvc := search.NewService(database, projSvc, tasksSvc, notesSvc)

	results, err := searchSvc.Search(ctx, user.ID, "Atlas")
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}

	if len(results) < 2 {
		t.Fatalf("expected at least 2 search results for 'Atlas', got: %d", len(results))
	}
}
