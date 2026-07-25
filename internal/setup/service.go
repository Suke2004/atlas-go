// Package setup implements the first-run onboarding wizard for Atlas.
package setup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Suke2004/atlas-go/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var ErrSetupAlreadyCompleted = errors.New("setup: first-run setup has already been completed")

// CreateFirstUserInput contains the form parameters submitted during first-run setup.
type CreateFirstUserInput struct {
	Username    string
	DisplayName string
	Email       string
	Password    string
	Timezone    string
	Theme       string
}

// Service provides business logic for the first-run onboarding flow.
type Service struct {
	db *db.DB
}

// NewService constructs a new setup Service.
func NewService(database *db.DB) *Service {
	return &Service{db: database}
}

// IsFirstRun checks whether any user exists in the database.
// If count == 0, returns true (first-run wizard needed).
func (s *Service) IsFirstRun(ctx context.Context) (bool, error) {
	count, err := s.db.CountUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("setup: failed to count users: %w", err)
	}
	return count == 0, nil
}

// CreateFirstUser creates the primary owner account and applies initial settings.
// Returns ErrSetupAlreadyCompleted if a user already exists.
func (s *Service) CreateFirstUser(ctx context.Context, input CreateFirstUserInput) (*db.User, error) {
	isFirst, err := s.IsFirstRun(ctx)
	if err != nil {
		return nil, err
	}
	if !isFirst {
		return nil, ErrSetupAlreadyCompleted
	}

	if input.Username == "" || input.Password == "" || input.DisplayName == "" {
		return nil, errors.New("setup: username, password, and display name are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("setup: failed to hash password: %w", err)
	}

	tz := input.Timezone
	if tz == "" {
		tz = "UTC"
	}

	var emailVal interface{}
	if input.Email != "" {
		emailVal = input.Email
	}

	user, err := s.db.CreateUser(ctx, db.CreateUserParams{
		Username:     input.Username,
		DisplayName:  input.DisplayName,
		Email:        emailVal,
		PasswordHash: string(hash),
		Timezone:     tz,
	})
	if err != nil {
		return nil, fmt.Errorf("setup: failed to create first user: %w", err)
	}

	// Persist initial theme setting
	theme := input.Theme
	if theme == "" {
		theme = "system"
	}
	_ = s.db.SetSetting(ctx, db.SetSettingParams{
		UserID: user.ID,
		Key:    "theme",
		Value:  theme,
	})

	return &user, nil
}

// SeedDemoData seeds initial sample projects, tasks, notes, and journal entries.
func (s *Service) SeedDemoData(ctx context.Context, userID int64) error {
	// 1. Demo Project
	proj, err := s.db.CreateProject(ctx, db.CreateProjectParams{
		UserID:      userID,
		Name:        "Atlas Operating System",
		Description: "Self-hosted personal OS built with Go & HTMX.",
		Status:      "active",
		Color:       "#3b82f6",
		TargetDate:  time.Now().AddDate(0, 1, 0),
	})
	if err != nil {
		return fmt.Errorf("setup: seed project: %w", err)
	}

	_ = s.db.UpdateProjectProgress(ctx, db.UpdateProjectProgressParams{
		ProgressPercentage: 35,
		ID:                 proj.ID,
		UserID:             userID,
	})

	_, _ = s.db.CreateMilestone(ctx, db.CreateMilestoneParams{
		ProjectID: proj.ID,
		Title:     "Phase 0+1: Architecture & Database Foundation",
		DueDate:   time.Now(),
	})

	// 2. Demo Tasks
	t1, err := s.db.CreateTask(ctx, db.CreateTaskParams{
		UserID:           userID,
		ProjectID:        proj.ID,
		Title:            "Explore Atlas Features",
		Description:      "Try out Projects, Tasks, Knowledge Base, and Journal entries.",
		Status:           "in_progress",
		Priority:         "high",
		EnergyLevel:      "medium",
		DueDate:          time.Now().AddDate(0, 0, 1),
		EstimatedMinutes: 30,
	})
	if err == nil {
		_ = s.db.AddTaskLabel(ctx, db.AddTaskLabelParams{
			TaskID: t1.ID,
			Label:  "onboarding",
		})
	}

	// 3. Demo Note
	note, err := s.db.CreateNote(ctx, db.CreateNoteParams{
		UserID:    userID,
		ProjectID: proj.ID,
		Title:     "Welcome to Atlas",
		Content:   "# Welcome to Atlas\n\nAtlas is your personal operating system. Everything is local, fast, and unified.",
		IsPinned:  true,
	})
	if err == nil {
		_ = s.db.AddNoteTag(ctx, db.AddNoteTagParams{
			NoteID: note.ID,
			Tag:    "welcome",
		})
	}

	// 4. Demo Journal Entry
	jEntry, err := s.db.UpsertJournalEntry(ctx, db.UpsertJournalEntryParams{
		UserID:       userID,
		EntryDate:    time.Now().Format("2006-01-02"),
		MoodRating:   5,
		EnergyRating: 5,
		SleepHours:  8.0,
		Summary:      "Completed Atlas setup wizard and initialized my personal OS!",
	})
	if err == nil {
		_, _ = s.db.AddJournalItem(ctx, db.AddJournalItemParams{
			EntryID:  jEntry.ID,
			Category: "win",
			Content:  "Configured my personal Atlas instance successfully.",
		})
	}

	return nil
}
