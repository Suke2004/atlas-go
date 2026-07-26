package notes

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	CreateNote(ctx context.Context, arg db.CreateNoteParams) (db.Note, error)
	GetNoteByID(ctx context.Context, arg db.GetNoteByIDParams) (db.Note, error)
	ListNotes(ctx context.Context, userID int64) ([]db.Note, error)
	UpdateNote(ctx context.Context, arg db.UpdateNoteParams) (db.Note, error)
	DeleteNote(ctx context.Context, arg db.DeleteNoteParams) error

	AddNoteTag(ctx context.Context, arg db.AddNoteTagParams) error
	ListNoteTags(ctx context.Context, noteID int64) ([]string, error)

	AddNoteLink(ctx context.Context, arg db.AddNoteLinkParams) error
	GetNoteBacklinks(ctx context.Context, arg db.GetNoteBacklinksParams) ([]db.Note, error)
}

type repository struct {
	database *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) CreateNote(ctx context.Context, arg db.CreateNoteParams) (db.Note, error) {
	return r.database.Queries.CreateNote(ctx, arg)
}

func (r *repository) GetNoteByID(ctx context.Context, arg db.GetNoteByIDParams) (db.Note, error) {
	return r.database.Queries.GetNoteByID(ctx, arg)
}

func (r *repository) ListNotes(ctx context.Context, userID int64) ([]db.Note, error) {
	return r.database.Queries.ListNotes(ctx, userID)
}

func (r *repository) UpdateNote(ctx context.Context, arg db.UpdateNoteParams) (db.Note, error) {
	return r.database.Queries.UpdateNote(ctx, arg)
}

func (r *repository) DeleteNote(ctx context.Context, arg db.DeleteNoteParams) error {
	return r.database.Queries.DeleteNote(ctx, arg)
}

func (r *repository) AddNoteTag(ctx context.Context, arg db.AddNoteTagParams) error {
	return r.database.Queries.AddNoteTag(ctx, arg)
}

func (r *repository) ListNoteTags(ctx context.Context, noteID int64) ([]string, error) {
	return r.database.Queries.ListNoteTags(ctx, noteID)
}

func (r *repository) AddNoteLink(ctx context.Context, arg db.AddNoteLinkParams) error {
	return r.database.Queries.AddNoteLink(ctx, arg)
}

func (r *repository) GetNoteBacklinks(ctx context.Context, arg db.GetNoteBacklinksParams) ([]db.Note, error) {
	return r.database.Queries.GetNoteBacklinks(ctx, arg)
}
