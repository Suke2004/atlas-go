package analytics

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

// Repository defines data access for analytics.
type Repository interface {
	GetDailyTaskCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetDailyNoteCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetDailyJournalCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetDailyLearningCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetDailyTransactionCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetDailyDocumentCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error)
	GetCategoryExpensesCurrentMonth(ctx context.Context, userID int64) ([]db.CategoryExpenseRow, error)
	GetMoodEnergyTrends30Days(ctx context.Context, userID int64) ([]db.MoodEnergyRow, error)
}

type repository struct {
	database *db.DB
}

// NewRepository constructs an analytics repository.
func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) GetDailyTaskCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyTaskCounts(ctx, userID, since)
}

func (r *repository) GetDailyNoteCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyNoteCounts(ctx, userID, since)
}

func (r *repository) GetDailyJournalCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyJournalCounts(ctx, userID, since)
}

func (r *repository) GetDailyLearningCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyLearningCounts(ctx, userID, since)
}

func (r *repository) GetDailyTransactionCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyTransactionCounts(ctx, userID, since)
}

func (r *repository) GetDailyDocumentCounts(ctx context.Context, userID int64, since string) ([]db.DailyCountRow, error) {
	return r.database.Queries.GetDailyDocumentCounts(ctx, userID, since)
}

func (r *repository) GetCategoryExpensesCurrentMonth(ctx context.Context, userID int64) ([]db.CategoryExpenseRow, error) {
	return r.database.Queries.GetCategoryExpensesCurrentMonth(ctx, userID)
}

func (r *repository) GetMoodEnergyTrends30Days(ctx context.Context, userID int64) ([]db.MoodEnergyRow, error) {
	return r.database.Queries.GetMoodEnergyTrends30Days(ctx, userID)
}
