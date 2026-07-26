package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/projects"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestProjects_ParseGitHubURL(t *testing.T) {
	tests := []struct {
		input     string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{"https://github.com/Suke2004/atlas-go", "Suke2004", "atlas-go", false},
		{"https://github.com/Suke2004/atlas-go.git", "Suke2004", "atlas-go", false},
		{"Suke2004/atlas-go", "Suke2004", "atlas-go", false},
		{"invalid-url", "", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		owner, repo, err := projects.ParseGitHubURL(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseGitHubURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if owner != tt.wantOwner || repo != tt.wantRepo {
			t.Errorf("ParseGitHubURL(%q) = (%q, %q), want (%q, %q)", tt.input, owner, repo, tt.wantOwner, tt.wantRepo)
		}
	}
}

func TestProjects_RecalculateProgress(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "projects_test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	if err := db.MigrateUp(database.Raw); err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	repo := projects.NewRepository(database)
	svc := projects.NewService(repo, nil)
	setupSvc := setup.NewService(database)
	ctx := context.Background()

	// 1. Create User
	user, err := setupSvc.CreateFirstUser(ctx, setup.CreateFirstUserInput{
		Username:    "projectuser",
		DisplayName: "Project Tester",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// 2. Create Project
	proj, err := svc.CreateProject(ctx, user.ID, projects.ProjectInput{
		Name:      "Atlas Engine",
		TechStack: "Go, SQLite, HTMX",
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	if proj.ProgressPercentage != 0 {
		t.Errorf("expected initial progress to be 0%%, got: %d%%", proj.ProgressPercentage)
	}

	// 3. Add 4 Milestones
	m1, err := svc.CreateMilestone(ctx, user.ID, proj.ID, "Phase 1 - DB", "")
	if err != nil {
		t.Fatalf("failed to create milestone: %v", err)
	}
	_, _ = svc.CreateMilestone(ctx, user.ID, proj.ID, "Phase 2 - Auth", "")
	_, _ = svc.CreateMilestone(ctx, user.ID, proj.ID, "Phase 3 - Layout", "")
	_, _ = svc.CreateMilestone(ctx, user.ID, proj.ID, "Phase 4 - Projects", "")

	// 4. Toggle 1 milestone -> expect 25% progress
	_, err = svc.ToggleMilestone(ctx, user.ID, proj.ID, m1.ID, true)
	if err != nil {
		t.Fatalf("failed to toggle milestone: %v", err)
	}

	updated, err := svc.GetProject(ctx, user.ID, proj.ID)
	if err != nil {
		t.Fatalf("failed to get updated project: %v", err)
	}

	if updated.Project.ProgressPercentage != 25 {
		t.Errorf("expected progress percentage to be 25%%, got: %d%%", updated.Project.ProgressPercentage)
	}

	// 5. Test GetProjectsSummary and tag filtering
	summary, err := svc.GetProjectsSummary(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get projects summary: %v", err)
	}

	if summary.TotalProjects != 1 || summary.ActiveProjects != 1 {
		t.Errorf("expected 1 total & active project, got total: %d, active: %d", summary.TotalProjects, summary.ActiveProjects)
	}

	filteredGo, err := svc.ListProjects(ctx, user.ID, "all", "Go", "")
	if err != nil || len(filteredGo) != 1 {
		t.Fatalf("expected 1 project matching tag 'Go', got: %v", filteredGo)
	}

	filteredRust, err := svc.ListProjects(ctx, user.ID, "all", "Rust", "")
	if err != nil || len(filteredRust) != 0 {
		t.Fatalf("expected 0 projects matching tag 'Rust', got: %v", filteredRust)
	}
}
