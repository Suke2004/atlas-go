package finance

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error)
	ListTransactions(ctx context.Context, userID int64) ([]db.Transaction, error)
	DeleteTransaction(ctx context.Context, arg db.DeleteTransactionParams) error
	GetFinanceSummary(ctx context.Context, userID int64) (db.GetFinanceSummaryRow, error)
}

type repository struct {
	database *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) CreateTransaction(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error) {
	return r.database.Queries.CreateTransaction(ctx, arg)
}

func (r *repository) ListTransactions(ctx context.Context, userID int64) ([]db.Transaction, error) {
	return r.database.Queries.ListTransactions(ctx, userID)
}

func (r *repository) DeleteTransaction(ctx context.Context, arg db.DeleteTransactionParams) error {
	return r.database.Queries.DeleteTransaction(ctx, arg)
}

func (r *repository) GetFinanceSummary(ctx context.Context, userID int64) (db.GetFinanceSummaryRow, error) {
	return r.database.Queries.GetFinanceSummary(ctx, userID)
}
