package learning

import "github.com/Suke2004/atlas-go/internal/db"

type TrackWithSessions struct {
	Track        db.LearningTrack     `json:"track"`
	Sessions     []db.LearningSession `json:"sessions"`
	TotalMinutes int64                `json:"total_minutes"`
}

type LearningSummary struct {
	TotalTracks     int64 `json:"total_tracks"`
	TotalStudyHours int64 `json:"total_study_hours"`
	ActiveStreaks   int64 `json:"active_streaks"`
	MasteryXP       int64 `json:"mastery_xp"`
}
