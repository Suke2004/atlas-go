package documents

import (
	"context"
	"database/sql"

	"github.com/Suke2004/atlas-go/internal/db"
)

// Repository defines the database access interface for documents.
type Repository interface {
	Create(ctx context.Context, p db.CreateDocumentParams) (db.Document, error)
	Get(ctx context.Context, userID, id int64) (db.Document, error)
	List(ctx context.Context, userID int64) ([]db.Document, error)
	UpdateMeta(ctx context.Context, p db.UpdateDocumentMetaParams) (db.Document, error)
	UpdateContent(ctx context.Context, id, userID int64, text string) error
	UpdateSummary(ctx context.Context, id, userID int64, summary string) error
	Delete(ctx context.Context, userID, id int64) error
	Count(ctx context.Context, userID int64) (int64, error)
}

type repository struct {
	database *db.DB
}

// NewRepository constructs a documents repository backed by SQLite.
func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) Create(ctx context.Context, p db.CreateDocumentParams) (db.Document, error) {
	return r.database.Queries.CreateDocument(ctx, p)
}

func (r *repository) Get(ctx context.Context, userID, id int64) (db.Document, error) {
	return r.database.Queries.GetDocument(ctx, db.GetDocumentParams{ID: id, UserID: userID})
}

func (r *repository) List(ctx context.Context, userID int64) ([]db.Document, error) {
	return r.database.Queries.ListDocuments(ctx, userID)
}

func (r *repository) UpdateMeta(ctx context.Context, p db.UpdateDocumentMetaParams) (db.Document, error) {
	return r.database.Queries.UpdateDocumentMeta(ctx, p)
}

func (r *repository) UpdateContent(ctx context.Context, id, userID int64, text string) error {
	return r.database.Queries.UpdateDocumentContent(ctx, db.UpdateDocumentContentParams{
		ID:          id,
		UserID:      userID,
		ContentText: sql.NullString{String: text, Valid: text != ""},
	})
}

func (r *repository) UpdateSummary(ctx context.Context, id, userID int64, summary string) error {
	return r.database.Queries.UpdateDocumentSummary(ctx, db.UpdateDocumentSummaryParams{
		ID:      id,
		UserID:  userID,
		Summary: sql.NullString{String: summary, Valid: summary != ""},
	})
}

func (r *repository) Delete(ctx context.Context, userID, id int64) error {
	return r.database.Queries.DeleteDocument(ctx, db.DeleteDocumentParams{ID: id, UserID: userID})
}

func (r *repository) Count(ctx context.Context, userID int64) (int64, error) {
	return r.database.Queries.CountDocuments(ctx, userID)
}
