package journal

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	UpsertJournalEntry(ctx context.Context, arg db.UpsertJournalEntryParams) (db.JournalEntry, error)
	GetJournalEntryByDate(ctx context.Context, arg db.GetJournalEntryByDateParams) (db.JournalEntry, error)
	ListJournalEntries(ctx context.Context, userID int64) ([]db.JournalEntry, error)
	AddJournalItem(ctx context.Context, arg db.AddJournalItemParams) (db.JournalItem, error)
	ListJournalItems(ctx context.Context, entryID int64) ([]db.JournalItem, error)
	DeleteJournalItem(ctx context.Context, id int64) error
}

type repository struct {
	database *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) UpsertJournalEntry(ctx context.Context, arg db.UpsertJournalEntryParams) (db.JournalEntry, error) {
	return r.database.Queries.UpsertJournalEntry(ctx, arg)
}

func (r *repository) GetJournalEntryByDate(ctx context.Context, arg db.GetJournalEntryByDateParams) (db.JournalEntry, error) {
	return r.database.Queries.GetJournalEntryByDate(ctx, arg)
}

func (r *repository) ListJournalEntries(ctx context.Context, userID int64) ([]db.JournalEntry, error) {
	return r.database.Queries.ListJournalEntries(ctx, userID)
}

func (r *repository) AddJournalItem(ctx context.Context, arg db.AddJournalItemParams) (db.JournalItem, error) {
	return r.database.Queries.AddJournalItem(ctx, arg)
}

func (r *repository) ListJournalItems(ctx context.Context, entryID int64) ([]db.JournalItem, error) {
	return r.database.Queries.ListJournalItems(ctx, entryID)
}

func (r *repository) DeleteJournalItem(ctx context.Context, id int64) error {
	return r.database.Queries.DeleteJournalItem(ctx, id)
}
