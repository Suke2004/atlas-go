package projects

type ProjectsSummary struct {
	TotalProjects           int64    `json:"total_projects"`
	ActiveProjects          int64    `json:"active_projects"`
	CompletedProjects       int64    `json:"completed_projects"`
	TotalMilestones         int64    `json:"total_milestones"`
	CompletedMilestones      int64    `json:"completed_milestones"`
	MilestoneCompletionRate int64    `json:"milestone_completion_rate"`
	TotalStars              int64    `json:"total_stars"`
	TotalForks              int64    `json:"total_forks"`
	TopTechStack            []string `json:"top_tech_stack"`
}
