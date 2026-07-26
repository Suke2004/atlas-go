package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/setup"
	"github.com/Suke2004/atlas-go/internal/tasks"
)

func TestTasks_ServiceAndProjectProgress(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "tasks_test.db")

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
		Username:    "taskuser",
		DisplayName: "Task User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	projRepo := projects.NewRepository(database)
	projSvc := projects.NewService(projRepo, nil)

	// Create project
	proj, err := projSvc.CreateProject(ctx, user.ID, projects.ProjectInput{
		Name: "Task Engine Integration",
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	tasksRepo := tasks.NewRepository(database)
	tasksSvc := tasks.NewService(tasksRepo, projRepo, projSvc)

	// Create task linked to project
	task, err := tasksSvc.CreateTask(ctx, user.ID, tasks.TaskInput{
		ProjectID:   proj.ID,
		Title:       "Build Tasks Module",
		Priority:    "high",
		EnergyLevel: "high",
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if task.Status != "todo" {
		t.Errorf("expected initial status 'todo', got: %s", task.Status)
	}

	// Update status to 'done'
	updatedTask, err := tasksSvc.UpdateTaskStatus(ctx, user.ID, task.ID, "done")
	if err != nil {
		t.Fatalf("failed to update task status: %v", err)
	}

	if updatedTask.Status != "done" {
		t.Errorf("expected status 'done', got: %s", updatedTask.Status)
	}

	summary, err := tasksSvc.GetTasksSummary(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get tasks summary: %v", err)
	}

	if summary.TotalTasks != 1 || summary.DoneTasks != 1 {
		t.Errorf("expected 1 total & done task, got total: %d, done: %d", summary.TotalTasks, summary.DoneTasks)
	}
}
