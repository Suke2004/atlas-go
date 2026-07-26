package tasks

import "github.com/Suke2004/atlas-go/internal/db"

type TasksSummary struct {
	TotalTasks      int64 `json:"total_tasks"`
	TodoTasks       int64 `json:"todo_tasks"`
	InProgressTasks int64 `json:"in_progress_tasks"`
	DoneTasks       int64 `json:"done_tasks"`
	TodayFocusTasks int64 `json:"today_focus_tasks"`
	CompletionRate  int64 `json:"completion_rate"`
}

type TaskWithDetails struct {
	Task         db.Task   `json:"task"`
	ProjectName  string    `json:"project_name"`
	Labels       []string  `json:"labels"`
	Dependencies []db.Task `json:"dependencies"`
}
