package unit

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/finance"
	"github.com/Suke2004/atlas-go/internal/setup"
)

func TestFinance_ServiceAndSummary(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "finance_test.db")

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
		Username:    "financeuser",
		DisplayName: "Finance User",
		Password:    "password123",
		Timezone:    "UTC",
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	finRepo := finance.NewRepository(database)
	finSvc := finance.NewService(finRepo, nil)

	// Add Income
	_, err = finSvc.CreateTransaction(ctx, user.ID, finance.TransactionInput{
		Amount:          5000.00,
		Type:            "income",
		Category:        "Salary",
		Description:     "Monthly engineering payout",
		TransactionDate: "2026-07-26",
	})
	if err != nil {
		t.Fatalf("failed to create income: %v", err)
	}

	// Add Expense
	_, err = finSvc.CreateTransaction(ctx, user.ID, finance.TransactionInput{
		Amount:          150.00,
		Type:            "expense",
		Category:        "Hosting",
		Description:     "Hetzner Cloud VPS",
		TransactionDate: "2026-07-26",
	})
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}

	summary, err := finSvc.GetFinanceSummary(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get finance summary: %v", err)
	}

	if summary.TotalIncome != 5000.00 {
		t.Errorf("expected income 5000, got: %f", summary.TotalIncome)
	}
	if summary.TotalExpenses != 150.00 {
		t.Errorf("expected expenses 150, got: %f", summary.TotalExpenses)
	}
	if summary.NetSavings != 4850.00 {
		t.Errorf("expected net savings 4850, got: %f", summary.NetSavings)
	}
}
