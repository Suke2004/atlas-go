package journal

import "github.com/Suke2004/atlas-go/internal/db"

type MindSyncPulse struct {
	TasksCompletedToday int64    `json:"tasks_completed_today"`
	CompletedTaskTitles []string `json:"completed_task_titles"`
	ProjectsUpdated     []string `json:"projects_updated"`
	NotesCreatedToday   int64    `json:"notes_created_today"`
}

type JournalEntryDetails struct {
	Entry    db.JournalEntry  `json:"entry"`
	Items    []db.JournalItem `json:"items"`
	MindSync MindSyncPulse    `json:"mind_sync"`
	Wins     []db.JournalItem `json:"wins"`
	Problems []db.JournalItem `json:"problems"`
	Ideas    []db.JournalItem `json:"ideas"`
	Tomorrow []db.JournalItem `json:"tomorrow"`
}

type JournalSummary struct {
	TotalEntries    int64   `json:"total_entries"`
	AvgMoodRating   float64 `json:"avg_mood_rating"`
	AvgEnergyRating float64 `json:"avg_energy_rating"`
	StreakDays      int     `json:"streak_days"`
}
