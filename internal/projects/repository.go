package projects

import (
	"context"

	"github.com/Suke2004/atlas-go/internal/db"
)

type Repository interface {
	CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error)
	GetProjectByID(ctx context.Context, arg db.GetProjectByIDParams) (db.Project, error)
	ListProjects(ctx context.Context, userID int64) ([]db.Project, error)
	UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error)
	UpdateGitHubStats(ctx context.Context, arg db.UpdateGitHubStatsParams) (db.Project, error)
	UpdateProjectProgress(ctx context.Context, arg db.UpdateProjectProgressParams) error
	DeleteProject(ctx context.Context, arg db.DeleteProjectParams) error

	CreateMilestone(ctx context.Context, arg db.CreateMilestoneParams) (db.Milestone, error)
	ListMilestonesByProject(ctx context.Context, projectID int64) ([]db.Milestone, error)
	ToggleMilestone(ctx context.Context, arg db.ToggleMilestoneParams) (db.Milestone, error)
	DeleteMilestone(ctx context.Context, arg db.DeleteMilestoneParams) error
}

type repository struct {
	queries *db.Queries
}

func NewRepository(database *db.DB) Repository {
	return &repository{
		queries: database.Queries,
	}
}

func (r *repository) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	return r.queries.CreateProject(ctx, arg)
}

func (r *repository) GetProjectByID(ctx context.Context, arg db.GetProjectByIDParams) (db.Project, error) {
	return r.queries.GetProjectByID(ctx, arg)
}

func (r *repository) ListProjects(ctx context.Context, userID int64) ([]db.Project, error) {
	return r.queries.ListProjects(ctx, userID)
}

func (r *repository) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	return r.queries.UpdateProject(ctx, arg)
}

func (r *repository) UpdateGitHubStats(ctx context.Context, arg db.UpdateGitHubStatsParams) (db.Project, error) {
	return r.queries.UpdateGitHubStats(ctx, arg)
}

func (r *repository) UpdateProjectProgress(ctx context.Context, arg db.UpdateProjectProgressParams) error {
	return r.queries.UpdateProjectProgress(ctx, arg)
}

func (r *repository) DeleteProject(ctx context.Context, arg db.DeleteProjectParams) error {
	return r.queries.DeleteProject(ctx, arg)
}

func (r *repository) CreateMilestone(ctx context.Context, arg db.CreateMilestoneParams) (db.Milestone, error) {
	return r.queries.CreateMilestone(ctx, arg)
}

func (r *repository) ListMilestonesByProject(ctx context.Context, projectID int64) ([]db.Milestone, error) {
	return r.queries.ListMilestonesByProject(ctx, projectID)
}

func (r *repository) ToggleMilestone(ctx context.Context, arg db.ToggleMilestoneParams) (db.Milestone, error) {
	return r.queries.ToggleMilestone(ctx, arg)
}

func (r *repository) DeleteMilestone(ctx context.Context, arg db.DeleteMilestoneParams) error {
	return r.queries.DeleteMilestone(ctx, arg)
}
