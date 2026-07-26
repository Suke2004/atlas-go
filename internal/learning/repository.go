package learning

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	CreateLearningTrack(ctx context.Context, arg db.CreateLearningTrackParams) (db.LearningTrack, error)
	ListLearningTracks(ctx context.Context, userID int64) ([]db.LearningTrack, error)
	DeleteLearningTrack(ctx context.Context, arg db.DeleteLearningTrackParams) error
	AddLearningSession(ctx context.Context, arg db.AddLearningSessionParams) (db.LearningSession, error)
	ListLearningSessions(ctx context.Context, trackID int64) ([]db.LearningSession, error)
}

type repository struct {
	database *db.DB
}

func NewRepository(database *db.DB) Repository {
	return &repository{database: database}
}

func (r *repository) CreateLearningTrack(ctx context.Context, arg db.CreateLearningTrackParams) (db.LearningTrack, error) {
	return r.database.Queries.CreateLearningTrack(ctx, arg)
}

func (r *repository) ListLearningTracks(ctx context.Context, userID int64) ([]db.LearningTrack, error) {
	return r.database.Queries.ListLearningTracks(ctx, userID)
}

func (r *repository) DeleteLearningTrack(ctx context.Context, arg db.DeleteLearningTrackParams) error {
	return r.database.Queries.DeleteLearningTrack(ctx, arg)
}

func (r *repository) AddLearningSession(ctx context.Context, arg db.AddLearningSessionParams) (db.LearningSession, error) {
	return r.database.Queries.AddLearningSession(ctx, arg)
}

func (r *repository) ListLearningSessions(ctx context.Context, trackID int64) ([]db.LearningSession, error) {
	return r.database.Queries.ListLearningSessions(ctx, trackID)
}
