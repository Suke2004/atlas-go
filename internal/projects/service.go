package projects

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Suke2004/atlas-go/internal/db"
	projecttemplates "github.com/Suke2004/atlas-go/web/templates/projects"
)

type ProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Color       string `json:"color"`
	TargetDate  string `json:"target_date"`
	GithubURL   string `json:"github_url"`
	TechStack   string `json:"tech_stack"`
}

type ProjectWithMilestones struct {
	Project    db.Project     `json:"project"`
	Milestones []db.Milestone `json:"milestones"`
}

type ProjectsSummary = projecttemplates.ProjectsSummary

type Service interface {
	CreateProject(ctx context.Context, userID int64, input ProjectInput) (db.Project, error)
	GetProject(ctx context.Context, userID, projectID int64) (*ProjectWithMilestones, error)
	ListProjects(ctx context.Context, userID int64, statusFilter string, tagFilter string, searchQuery string) ([]db.Project, error)
	GetProjectsSummary(ctx context.Context, userID int64) (ProjectsSummary, error)
	UpdateProject(ctx context.Context, userID, projectID int64, input ProjectInput) (db.Project, error)
	DeleteProject(ctx context.Context, userID, projectID int64) error
	SyncGitHubStats(ctx context.Context, userID, projectID int64) (db.Project, error)
	ImportGitHubIssues(ctx context.Context, userID, projectID int64) ([]db.Milestone, error)

	CreateMilestone(ctx context.Context, userID, projectID int64, title string, dueDate string) (db.Milestone, error)
	ToggleMilestone(ctx context.Context, userID, projectID, milestoneID int64, isCompleted bool) (db.Milestone, error)
	DeleteMilestone(ctx context.Context, userID, projectID, milestoneID int64) error
	RecalculateProgress(ctx context.Context, userID, projectID int64) error
}

type service struct {
	repo         Repository
	githubClient *GitHubClient
}

func NewService(repo Repository, ghClient *GitHubClient) Service {
	if ghClient == nil {
		ghClient = NewGitHubClient()
	}
	return &service{
		repo:         repo,
		githubClient: ghClient,
	}
}

func (s *service) CreateProject(ctx context.Context, userID int64, input ProjectInput) (db.Project, error) {
	if strings.TrimSpace(input.Name) == "" {
		return db.Project{}, fmt.Errorf("project name is required")
	}
	if input.Status == "" {
		input.Status = "active"
	}
	if input.Color == "" {
		input.Color = "#6366f1"
	}

	var ghStars, ghForks, ghIssues int64
	var ghLang string
	var ghLastPush interface{}
	techStack := strings.TrimSpace(input.TechStack)

	// Auto-fetch GitHub stats if GitHub URL provided
	if input.GithubURL != "" {
		owner, repo, err := ParseGitHubURL(input.GithubURL)
		if err == nil {
			stats, err := s.githubClient.FetchRepoStats(ctx, owner, repo)
			if err == nil {
				ghStars = stats.Stars
				ghForks = stats.Forks
				ghIssues = stats.OpenIssues
				ghLang = stats.PrimaryLanguage
				if !stats.LastPushedAt.IsZero() {
					ghLastPush = stats.LastPushedAt.Format("2006-01-02 15:04:05")
				}
				if techStack == "" && len(stats.TechStack) > 0 {
					techStack = strings.Join(stats.TechStack, ", ")
				}
			}
		}
	}

	var targetDate interface{}
	if strings.TrimSpace(input.TargetDate) != "" {
		targetDate = strings.TrimSpace(input.TargetDate)
	}

	arg := db.CreateProjectParams{
		UserID:             userID,
		Name:               strings.TrimSpace(input.Name),
		Description:        strings.TrimSpace(input.Description),
		Status:             input.Status,
		Color:              input.Color,
		TargetDate:         targetDate,
		GithubUrl:          strings.TrimSpace(input.GithubURL),
		GithubStars:        ghStars,
		GithubForks:        ghForks,
		GithubOpenIssues:   ghIssues,
		GithubLanguage:     ghLang,
		GithubLastPushedAt: ghLastPush,
		TechStack:          techStack,
	}

	return s.repo.CreateProject(ctx, arg)
}

func (s *service) GetProject(ctx context.Context, userID, projectID int64) (*ProjectWithMilestones, error) {
	project, err := s.repo.GetProjectByID(ctx, db.GetProjectByIDParams{ID: projectID, UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	milestones, err := s.repo.ListMilestonesByProject(ctx, projectID)
	if err != nil {
		milestones = []db.Milestone{}
	}

	return &ProjectWithMilestones{
		Project:    project,
		Milestones: milestones,
	}, nil
}

func (s *service) ListProjects(ctx context.Context, userID int64, statusFilter string, tagFilter string, searchQuery string) ([]db.Project, error) {
	allProjects, err := s.repo.ListProjects(ctx, userID)
	if err != nil {
		return nil, err
	}

	var filtered []db.Project
	tagLower := strings.ToLower(strings.TrimSpace(tagFilter))
	searchLower := strings.ToLower(strings.TrimSpace(searchQuery))

	for _, p := range allProjects {
		// 1. Status Filter
		if statusFilter != "" && statusFilter != "all" && !strings.EqualFold(p.Status, statusFilter) {
			continue
		}

		// 2. Tag Filter
		if tagLower != "" {
			matchedTag := false
			for _, tag := range strings.Split(p.TechStack, ",") {
				if strings.TrimSpace(strings.ToLower(tag)) == tagLower {
					matchedTag = true
					break
				}
			}
			if !matchedTag {
				continue
			}
		}

		// 3. Search Query
		if searchLower != "" {
			nameMatch := strings.Contains(strings.ToLower(p.Name), searchLower)
			descMatch := strings.Contains(strings.ToLower(p.Description), searchLower)
			repoMatch := strings.Contains(strings.ToLower(p.GithubUrl), searchLower)
			if !nameMatch && !descMatch && !repoMatch {
				continue
			}
		}

		filtered = append(filtered, p)
	}

	return filtered, nil
}

func (s *service) GetProjectsSummary(ctx context.Context, userID int64) (ProjectsSummary, error) {
	allProjects, err := s.repo.ListProjects(ctx, userID)
	if err != nil {
		return ProjectsSummary{}, err
	}

	var summary ProjectsSummary
	summary.TotalProjects = int64(len(allProjects))
	techMap := make(map[string]int)

	for _, p := range allProjects {
		if strings.EqualFold(p.Status, "active") {
			summary.ActiveProjects++
		} else if strings.EqualFold(p.Status, "completed") {
			summary.CompletedProjects++
		}

		summary.TotalStars += p.GithubStars
		summary.TotalForks += p.GithubForks

		// Milestones
		milestones, err := s.repo.ListMilestonesByProject(ctx, p.ID)
		if err == nil {
			summary.TotalMilestones += int64(len(milestones))
			for _, m := range milestones {
				if m.IsCompleted {
					summary.CompletedMilestones++
				}
			}
		}

		// Tech Stack breakdown
		if p.TechStack != "" {
			for _, tag := range strings.Split(p.TechStack, ",") {
				t := strings.TrimSpace(tag)
				if t != "" {
					techMap[t]++
				}
			}
		}
	}

	if summary.TotalMilestones > 0 {
		summary.MilestoneCompletionRate = (summary.CompletedMilestones * 100) / summary.TotalMilestones
	}

	// Extract top tech stack tags
	for tag := range techMap {
		summary.TopTechStack = append(summary.TopTechStack, tag)
		if len(summary.TopTechStack) >= 6 {
			break
		}
	}

	return summary, nil
}

func (s *service) UpdateProject(ctx context.Context, userID, projectID int64, input ProjectInput) (db.Project, error) {
	if strings.TrimSpace(input.Name) == "" {
		return db.Project{}, fmt.Errorf("project name is required")
	}

	var targetDate interface{}
	if strings.TrimSpace(input.TargetDate) != "" {
		targetDate = strings.TrimSpace(input.TargetDate)
	}

	arg := db.UpdateProjectParams{
		Name:        strings.TrimSpace(input.Name),
		Description: strings.TrimSpace(input.Description),
		Status:      input.Status,
		Color:       input.Color,
		TargetDate:  targetDate,
		GithubUrl:   strings.TrimSpace(input.GithubURL),
		TechStack:   strings.TrimSpace(input.TechStack),
		ID:          projectID,
		UserID:      userID,
	}

	return s.repo.UpdateProject(ctx, arg)
}

func (s *service) DeleteProject(ctx context.Context, userID, projectID int64) error {
	return s.repo.DeleteProject(ctx, db.DeleteProjectParams{ID: projectID, UserID: userID})
}

func (s *service) SyncGitHubStats(ctx context.Context, userID, projectID int64) (db.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, db.GetProjectByIDParams{ID: projectID, UserID: userID})
	if err != nil {
		return db.Project{}, err
	}

	if project.GithubUrl == "" {
		return project, fmt.Errorf("project does not have a github url configured")
	}

	owner, repo, err := ParseGitHubURL(project.GithubUrl)
	if err != nil {
		return project, fmt.Errorf("invalid github url: %w", err)
	}

	stats, err := s.githubClient.FetchRepoStats(ctx, owner, repo)
	if err != nil {
		return project, fmt.Errorf("failed to sync github stats: %w", err)
	}

	techStack := project.TechStack
	if len(stats.TechStack) > 0 {
		techStack = strings.Join(stats.TechStack, ", ")
	}

	var ghLastPush interface{}
	if !stats.LastPushedAt.IsZero() {
		ghLastPush = stats.LastPushedAt.Format("2006-01-02 15:04:05")
	}

	arg := db.UpdateGitHubStatsParams{
		GithubStars:        stats.Stars,
		GithubForks:        stats.Forks,
		GithubOpenIssues:   stats.OpenIssues,
		GithubLanguage:     stats.PrimaryLanguage,
		GithubLastPushedAt: ghLastPush,
		TechStack:          techStack,
		ID:                 projectID,
		UserID:             userID,
	}

	return s.repo.UpdateGitHubStats(ctx, arg)
}

func (s *service) ImportGitHubIssues(ctx context.Context, userID, projectID int64) ([]db.Milestone, error) {
	project, err := s.repo.GetProjectByID(ctx, db.GetProjectByIDParams{ID: projectID, UserID: userID})
	if err != nil {
		return nil, err
	}

	owner, repo, err := ParseGitHubURL(project.GithubUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid github url: %w", err)
	}

	issues, err := s.githubClient.FetchOpenIssues(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch github issues: %w", err)
	}

	var created []db.Milestone
	for _, issue := range issues {
		title := fmt.Sprintf("Issue: %s", issue.Title)
		m, err := s.repo.CreateMilestone(ctx, db.CreateMilestoneParams{
			ProjectID: projectID,
			Title:     title,
			DueDate:   sql.NullTime{},
		})
		if err == nil {
			created = append(created, m)
		}
	}

	_ = s.RecalculateProgress(ctx, userID, projectID)
	return created, nil
}

func (s *service) CreateMilestone(ctx context.Context, userID, projectID int64, title string, dueDateStr string) (db.Milestone, error) {
	if strings.TrimSpace(title) == "" {
		return db.Milestone{}, fmt.Errorf("milestone title is required")
	}

	var dueDate interface{}
	if strings.TrimSpace(dueDateStr) != "" {
		dueDate = strings.TrimSpace(dueDateStr)
	}

	m, err := s.repo.CreateMilestone(ctx, db.CreateMilestoneParams{
		ProjectID: projectID,
		Title:     strings.TrimSpace(title),
		DueDate:   dueDate,
	})
	if err != nil {
		return db.Milestone{}, err
	}

	_ = s.RecalculateProgress(ctx, userID, projectID)
	return m, nil
}

func (s *service) ToggleMilestone(ctx context.Context, userID, projectID, milestoneID int64, isCompleted bool) (db.Milestone, error) {
	arg := db.ToggleMilestoneParams{
		IsCompleted:   isCompleted,
		IsCompleted_2: isCompleted,
		ID:            milestoneID,
		ProjectID:     projectID,
	}

	m, err := s.repo.ToggleMilestone(ctx, arg)
	if err != nil {
		return db.Milestone{}, err
	}

	_ = s.RecalculateProgress(ctx, userID, projectID)
	return m, nil
}

func (s *service) DeleteMilestone(ctx context.Context, userID, projectID, milestoneID int64) error {
	err := s.repo.DeleteMilestone(ctx, db.DeleteMilestoneParams{ID: milestoneID, ProjectID: projectID})
	if err != nil {
		return err
	}

	_ = s.RecalculateProgress(ctx, userID, projectID)
	return nil
}

func (s *service) RecalculateProgress(ctx context.Context, userID, projectID int64) error {
	milestones, err := s.repo.ListMilestonesByProject(ctx, projectID)
	if err != nil || len(milestones) == 0 {
		return s.repo.UpdateProjectProgress(ctx, db.UpdateProjectProgressParams{
			ProgressPercentage: 0,
			ID:                 projectID,
			UserID:             userID,
		})
	}

	var completed int64
	for _, m := range milestones {
		if m.IsCompleted {
			completed++
		}
	}

	pct := (completed * 100) / int64(len(milestones))
	return s.repo.UpdateProjectProgress(ctx, db.UpdateProjectProgressParams{
		ProgressPercentage: pct,
		ID:                 projectID,
		UserID:             userID,
	})
}
