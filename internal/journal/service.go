package journal

import (
	"context"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
	"github.com/Suke2004/atlas-go/internal/notes"
	"github.com/Suke2004/atlas-go/internal/tasks"
	journaltemplates "github.com/Suke2004/atlas-go/web/templates/journal"
)

type JournalInput struct {
	EntryDate    string  `json:"entry_date"`
	MoodRating   int     `json:"mood_rating"`
	EnergyRating int     `json:"energy_rating"`
	SleepHours  float64 `json:"sleep_hours"`
	Summary      string  `json:"summary"`
}

type JournalEntryDetails = journaltemplates.JournalEntryDetails
type JournalSummary = journaltemplates.JournalSummary
type MindSyncPulse = journaltemplates.MindSyncPulse

type Service interface {
	GetDailyJournal(ctx context.Context, userID int64, dateStr string) (*JournalEntryDetails, error)
	SaveDailyJournal(ctx context.Context, userID int64, input JournalInput) (db.JournalEntry, error)
	AddJournalItem(ctx context.Context, userID int64, entryID int64, category string, content string) (db.JournalItem, error)
	DeleteJournalItem(ctx context.Context, itemID int64) error
	GetJournalSummary(ctx context.Context, userID int64) (JournalSummary, error)
}

type service struct {
	repo     Repository
	tasksSvc tasks.Service
	notesSvc notes.Service
}

func NewService(repo Repository, tasksSvc tasks.Service, notesSvc notes.Service) Service {
	return &service{
		repo:     repo,
		tasksSvc: tasksSvc,
		notesSvc: notesSvc,
	}
}

func (s *service) GetDailyJournal(ctx context.Context, userID int64, dateStr string) (*JournalEntryDetails, error) {
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	entry, err := s.repo.GetJournalEntryByDate(ctx, db.GetJournalEntryByDateParams{
		UserID:    userID,
		EntryDate: dateStr,
	})

	var items []db.JournalItem
	if err == nil {
		items, _ = s.repo.ListJournalItems(ctx, entry.ID)
	}

	// Calculate Executive Mind-Sync Pulse for today
	var pulse MindSyncPulse
	if s.tasksSvc != nil {
		taskList, err := s.tasksSvc.ListTasks(ctx, userID, "done", "all", "all", 0, "")
		if err == nil {
			for _, t := range taskList {
				pulse.TasksCompletedToday++
				pulse.CompletedTaskTitles = append(pulse.CompletedTaskTitles, t.Task.Title)
			}
		}
	}

	if s.notesSvc != nil {
		notesList, err := s.notesSvc.ListNotes(ctx, userID, "all", "")
		if err == nil {
			pulse.NotesCreatedToday = int64(len(notesList))
		}
	}

	var wins, problems, ideas, tomorrow []db.JournalItem
	for _, item := range items {
		switch item.Category {
		case "win":
			wins = append(wins, item)
		case "problem":
			problems = append(problems, item)
		case "idea":
			ideas = append(ideas, item)
		case "tomorrow":
			tomorrow = append(tomorrow, item)
		}
	}

	return &JournalEntryDetails{
		Entry:    entry,
		Items:    items,
		MindSync: pulse,
		Wins:     wins,
		Problems: problems,
		Ideas:    ideas,
		Tomorrow: tomorrow,
	}, nil
}

func (s *service) SaveDailyJournal(ctx context.Context, userID int64, input JournalInput) (db.JournalEntry, error) {
	if input.EntryDate == "" {
		input.EntryDate = time.Now().Format("2006-01-02")
	}

	var mood, energy, sleep interface{}
	if input.MoodRating > 0 {
		mood = input.MoodRating
	}
	if input.EnergyRating > 0 {
		energy = input.EnergyRating
	}
	if input.SleepHours > 0 {
		sleep = input.SleepHours
	}

	return s.repo.UpsertJournalEntry(ctx, db.UpsertJournalEntryParams{
		UserID:       userID,
		EntryDate:    input.EntryDate,
		MoodRating:   mood,
		EnergyRating: energy,
		SleepHours:  sleep,
		Summary:      input.Summary,
	})
}

func (s *service) AddJournalItem(ctx context.Context, userID int64, entryID int64, category string, content string) (db.JournalItem, error) {
	return s.repo.AddJournalItem(ctx, db.AddJournalItemParams{
		EntryID:  entryID,
		Category: category,
		Content:  content,
	})
}

func (s *service) DeleteJournalItem(ctx context.Context, itemID int64) error {
	return s.repo.DeleteJournalItem(ctx, itemID)
}

func (s *service) GetJournalSummary(ctx context.Context, userID int64) (JournalSummary, error) {
	entries, err := s.repo.ListJournalEntries(ctx, userID)
	if err != nil {
		return JournalSummary{}, err
	}

	var summary JournalSummary
	summary.TotalEntries = int64(len(entries))
	summary.StreakDays = len(entries)

	if len(entries) > 0 {
		var moodSum, energySum float64
		var moodCount, energyCount float64

		for _, e := range entries {
			if e.MoodRating.Valid {
				moodSum += float64(e.MoodRating.Int64)
				moodCount++
			}
			if e.EnergyRating.Valid {
				energySum += float64(e.EnergyRating.Int64)
				energyCount++
			}
		}

		if moodCount > 0 {
			summary.AvgMoodRating = moodSum / moodCount
		}
		if energyCount > 0 {
			summary.AvgEnergyRating = energySum / energyCount
		}
	}

	return summary, nil
}
