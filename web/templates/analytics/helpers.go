package analytics

import (
	"github.com/Suke2004/atlas-go/internal/db"
)

// DailyContribution represents activity count on a specific date.
type DailyContribution struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
	Level int    `json:"level"`
}

// SummaryMetrics aggregates overall activity metrics.
type SummaryMetrics struct {
	TotalContributions int                     `json:"total_contributions"`
	CurrentStreak      int                     `json:"current_streak"`
	LongestStreak      int                     `json:"longest_streak"`
	CategoryExpenses   []db.CategoryExpenseRow `json:"category_expenses"`
	MoodEnergyTrends   []db.MoodEnergyRow      `json:"mood_energy_trends"`
}

// AnalyticsData contains complete dataset for rendering the analytics dashboard.
type AnalyticsData struct {
	Heatmap      []DailyContribution `json:"heatmap"`
	Metrics      SummaryMetrics      `json:"metrics"`
	WeeklyReview string              `json:"weekly_review"`
}
