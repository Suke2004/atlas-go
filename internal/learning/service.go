package learning

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
	learningtemplates "github.com/Suke2004/atlas-go/web/templates/learning"
)

type TrackInput struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type TrackWithSessions = learningtemplates.TrackWithSessions
type LearningSummary = learningtemplates.LearningSummary

type Service interface {
	CreateTrack(ctx context.Context, userID int64, input TrackInput) (db.LearningTrack, error)
	ListTracks(ctx context.Context, userID int64) ([]TrackWithSessions, error)
	AddSession(ctx context.Context, userID int64, trackID int64, minutes int64, summary string) (db.LearningSession, error)
	GetLearningSummary(ctx context.Context, userID int64) (LearningSummary, error)
	DeleteTrack(ctx context.Context, userID, trackID int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateTrack(ctx context.Context, userID int64, input TrackInput) (db.LearningTrack, error) {
	if strings.TrimSpace(input.Title) == "" {
		return db.LearningTrack{}, fmt.Errorf("title is required")
	}

	return s.repo.CreateLearningTrack(ctx, db.CreateLearningTrackParams{
		UserID:        userID,
		Title:         strings.TrimSpace(input.Title),
		Category:      input.Category,
		Description:   strings.TrimSpace(input.Description),
		CurrentStreak: 1,
		LongestStreak: 1,
	})
}

func (s *service) ListTracks(ctx context.Context, userID int64) ([]TrackWithSessions, error) {
	tracks, err := s.repo.ListLearningTracks(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []TrackWithSessions
	for _, t := range tracks {
		sessions, _ := s.repo.ListLearningSessions(ctx, t.ID)
		var totalMins int64
		for _, sess := range sessions {
			totalMins += sess.DurationMinutes
		}
		result = append(result, TrackWithSessions{
			Track:        t,
			Sessions:     sessions,
			TotalMinutes: totalMins,
		})
	}

	return result, nil
}

func (s *service) AddSession(ctx context.Context, userID int64, trackID int64, minutes int64, summary string) (db.LearningSession, error) {
	if minutes <= 0 {
		minutes = 30
	}

	return s.repo.AddLearningSession(ctx, db.AddLearningSessionParams{
		TrackID:         trackID,
		DurationMinutes: minutes,
		Summary:         summary,
		SessionDate:     time.Now().Format("2006-01-02"),
	})
}

func (s *service) GetLearningSummary(ctx context.Context, userID int64) (LearningSummary, error) {
	tracks, err := s.repo.ListLearningTracks(ctx, userID)
	if err != nil {
		return LearningSummary{}, err
	}

	var summary LearningSummary
	summary.TotalTracks = int64(len(tracks))

	var totalMins int64
	for _, t := range tracks {
		if t.CurrentStreak > 0 {
			summary.ActiveStreaks++
		}
		sessions, _ := s.repo.ListLearningSessions(ctx, t.ID)
		for _, sess := range sessions {
			totalMins += sess.DurationMinutes
		}
	}

	summary.TotalStudyHours = totalMins / 60
	summary.MasteryXP = (summary.TotalStudyHours * 100) + (summary.ActiveStreaks * 50)

	return summary, nil
}

func (s *service) DeleteTrack(ctx context.Context, userID, trackID int64) error {
	return s.repo.DeleteLearningTrack(ctx, db.DeleteLearningTrackParams{ID: trackID, UserID: userID})
}
