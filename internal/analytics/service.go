package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/Suke2004/atlas-go/internal/ai"
	"github.com/Suke2004/atlas-go/internal/db"
	analyticstemplates "github.com/Suke2004/atlas-go/web/templates/analytics"
)

// Service defines analytics business logic.
type Service interface {
	GetAnalyticsData(ctx context.Context, userID int64) (analyticstemplates.AnalyticsData, error)
	GenerateWeeklyReview(ctx context.Context, userID int64) (string, error)
}

type service struct {
	repo       Repository
	aiProvider ai.Provider
}

// NewService constructs an analytics service.
func NewService(repo Repository, aiProvider ai.Provider) Service {
	return &service{
		repo:       repo,
		aiProvider: aiProvider,
	}
}

func (s *service) GetAnalyticsData(ctx context.Context, userID int64) (analyticstemplates.AnalyticsData, error) {
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0) // 365 days ago
	sinceStr := startDate.Format("2006-01-02")

	activityMap := make(map[string]int)

	// Helper to merge counts
	addCounts := func(rows []db.DailyCountRow) {
		for _, r := range rows {
			activityMap[r.DateStr] += int(r.Count)
		}
	}

	if tasks, err := s.repo.GetDailyTaskCounts(ctx, userID, sinceStr); err == nil {
		addCounts(tasks)
	}
	if notes, err := s.repo.GetDailyNoteCounts(ctx, userID, sinceStr); err == nil {
		addCounts(notes)
	}
	if journal, err := s.repo.GetDailyJournalCounts(ctx, userID, sinceStr); err == nil {
		addCounts(journal)
	}
	if learning, err := s.repo.GetDailyLearningCounts(ctx, userID, sinceStr); err == nil {
		addCounts(learning)
	}
	if transactions, err := s.repo.GetDailyTransactionCounts(ctx, userID, sinceStr); err == nil {
		addCounts(transactions)
	}
	if docs, err := s.repo.GetDailyDocumentCounts(ctx, userID, sinceStr); err == nil {
		addCounts(docs)
	}

	// Build 365 days slice
	heatmap := make([]analyticstemplates.DailyContribution, 0, 365)
	totalContribs := 0
	currStreak := 0
	maxStreak := 0
	tempStreak := 0

	for d := startDate; !d.After(now); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		cnt := activityMap[dateStr]
		totalContribs += cnt

		level := 0
		if cnt > 0 {
			switch {
			case cnt >= 8:
				level = 4
			case cnt >= 5:
				level = 3
			case cnt >= 3:
				level = 2
			default:
				level = 1
			}
			tempStreak++
			if tempStreak > maxStreak {
				maxStreak = tempStreak
			}
		} else {
			tempStreak = 0
		}

		heatmap = append(heatmap, analyticstemplates.DailyContribution{
			Date:  dateStr,
			Count: cnt,
			Level: level,
		})
	}

	// Calculate current streak backwards from today
	for i := len(heatmap) - 1; i >= 0; i-- {
		if heatmap[i].Count > 0 {
			currStreak++
		} else if i == len(heatmap)-1 {
			// Today might not have activity yet, check yesterday
			continue
		} else {
			break
		}
	}

	expenses, _ := s.repo.GetCategoryExpensesCurrentMonth(ctx, userID)
	trends, _ := s.repo.GetMoodEnergyTrends30Days(ctx, userID)

	return analyticstemplates.AnalyticsData{
		Heatmap: heatmap,
		Metrics: analyticstemplates.SummaryMetrics{
			TotalContributions: totalContribs,
			CurrentStreak:      currStreak,
			LongestStreak:      maxStreak,
			CategoryExpenses:   expenses,
			MoodEnergyTrends:   trends,
		},
	}, nil
}

func (s *service) GenerateWeeklyReview(ctx context.Context, userID int64) (string, error) {
	if s.aiProvider == nil {
		return "", fmt.Errorf("AI provider not configured — set one in Settings")
	}

	data, err := s.GetAnalyticsData(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to gather analytics for AI review: %w", err)
	}

	prompt := fmt.Sprintf(
		"Synthesize a concise executive weekly productivity review for the user based on these 365-day telemetry numbers:\n"+
			"- Total Contributions Across System: %d\n"+
			"- Current Activity Streak: %d days\n"+
			"- Longest Activity Streak: %d days\n"+
			"- Active Mood/Energy Trend Points: %d entries\n"+
			"Format in clear markdown sections: '🚀 Weekly Accomplishments', '📊 Well-Being & Energy Insight', and '🎯 Actionable Focus for Next Week'. Keep it inspiring and under 250 words.",
		data.Metrics.TotalContributions,
		data.Metrics.CurrentStreak,
		data.Metrics.LongestStreak,
		len(data.Metrics.MoodEnergyTrends),
	)

	messages := []ai.Message{
		{Role: ai.RoleSystem, Content: "You are Atlas Analytics AI, an executive coach analyzing developer telemetry data."},
		{Role: ai.RoleUser, Content: prompt},
	}

	return s.aiProvider.Complete(ctx, messages)
}
