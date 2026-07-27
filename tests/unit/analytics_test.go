package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/analytics"
	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestAnalytics_ServiceAndDataAggregation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "analytics_test.db")

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
		Username:    "analyticsuser",
		DisplayName: "Analytics User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	repo := analytics.NewRepository(database)
	svc := analytics.NewService(repo, nil)

	data, err := svc.GetAnalyticsData(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get analytics data: %v", err)
	}

	if len(data.Heatmap) == 0 {
		t.Errorf("expected non-empty heatmap array, got 0")
	}

	if data.Metrics.TotalContributions < 0 {
		t.Errorf("invalid total contributions: %d", data.Metrics.TotalContributions)
	}
}
