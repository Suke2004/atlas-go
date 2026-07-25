// Package main implements the development seed script to populate Atlas with initial demo data.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Suke2004/atlas-go/internal/config"
	"github.com/Suke2004/atlas-go/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	// Check if user already exists
	count, err := database.CountUsers(ctx)
	if err != nil {
		log.Fatalf("Failed to count users: %v", err)
	}

	if count > 0 {
		fmt.Println("Database already contains user data. Skipping seed.")
		os.Exit(0)
	}

	// 1. Create default demo user
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user, err := database.CreateUser(ctx, db.CreateUserParams{
		Username:     "atlas_user",
		DisplayName:  "Atlas Explorer",
		Email:        sqlNullString("user@atlas.local"),
		PasswordHash: string(hash),
		Timezone:     "UTC",
	})
	if err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}

	fmt.Printf("Created demo user: %s (ID: %d)\n", user.Username, user.ID)

	// 2. Create demo projects
	atlasProj, err := database.CreateProject(ctx, db.CreateProjectParams{
		UserID:             user.ID,
		Name:               "Atlas Operating System",
		Description:        "Self-hosted personal OS built with Go and HTMX.",
		Status:             "active",
		Color:              "#3b82f6",
		TargetDate:         sqlNullTime(time.Now().AddDate(0, 1, 0)),
	})
	if err == nil {
		_ = database.UpdateProjectProgress(ctx, db.UpdateProjectProgressParams{
			ProgressPercentage: 35,
			ID:                 atlasProj.ID,
			UserID:             user.ID,
		})

		_, _ = database.CreateMilestone(ctx, db.CreateMilestoneParams{
			ProjectID: atlasProj.ID,
			Title:     "Phase 0: Scaffold & Architecture",
			DueDate:   sqlNullTime(time.Now()),
		})
		_, _ = database.CreateMilestone(ctx, db.CreateMilestoneParams{
			ProjectID: atlasProj.ID,
			Title:     "Phase 1: Database Foundation & Migrations",
			DueDate:   sqlNullTime(time.Now().AddDate(0, 0, 2)),
		})
	}

	// 3. Create demo tasks
	task1, err := database.CreateTask(ctx, db.CreateTaskParams{
		UserID:           user.ID,
		ProjectID:        sqlNullInt64(atlasProj.ID),
		Title:            "Complete Phase 1 Migrations & sqlc generation",
		Description:      "Write Goose SQL migration files 001-009 and compile sqlc query handlers.",
		Status:           "in_progress",
		Priority:         "high",
		EnergyLevel:      "high",
		DueDate:          sqlNullTime(time.Now().AddDate(0, 0, 1)),
		EstimatedMinutes: sqlNullInt64(120),
	})

	if err == nil {
		_ = database.AddTaskLabel(ctx, db.AddTaskLabelParams{
			TaskID: task1.ID,
			Label:  "database",
		})
	}

	_, _ = database.CreateTask(ctx, db.CreateTaskParams{
		UserID:           user.ID,
		ProjectID:        sqlNullInt64(atlasProj.ID),
		Title:            "Implement First-Run Setup Wizard & Auth",
		Description:      "Build /setup wizard flow and session-based login middleware.",
		Status:           "todo",
		Priority:         "critical",
		EnergyLevel:      "medium",
		DueDate:          sqlNullTime(time.Now().AddDate(0, 0, 3)),
		EstimatedMinutes: sqlNullInt64(180),
	})

	// 4. Create demo Knowledge Base notes
	note, err := database.CreateNote(ctx, db.CreateNoteParams{
		UserID:    user.ID,
		ProjectID: sqlNullInt64(atlasProj.ID),
		Title:     "Atlas System Architecture & Philosophy",
		Content:   "# Atlas Architecture\n\nAtlas is a server-rendered personal operating system. Every interaction uses HTMX partial updates for smooth, low-latency performance without SPA complexity.",
		IsPinned:  true,
	})
	if err == nil {
		_ = database.AddNoteTag(ctx, db.AddNoteTagParams{
			NoteID: note.ID,
			Tag:    "architecture",
		})
	}

	// 5. Create demo Journal Entry
	todayStr := time.Now().Format("2006-01-02")
	jEntry, err := database.UpsertJournalEntry(ctx, db.UpsertJournalEntryParams{
		UserID:       user.ID,
		EntryDate:    todayStr,
		MoodRating:   sqlNullInt64(5),
		EnergyRating: sqlNullInt64(4),
		SleepHours:  sqlNullFloat64(7.5),
		Summary:      "Great progress on Phase 0 and Phase 1 architecture of Atlas today!",
	})
	if err == nil {
		_, _ = database.AddJournalItem(ctx, db.AddJournalItemParams{
			EntryID:  jEntry.ID,
			Category: "win",
			Content:  "Scaffolded Atlas repo with complete 135KB documentation & guidelines.",
		})
	}

	fmt.Println("Demo data seeded successfully!")
}

func sqlNullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func sqlNullInt64(i int64) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

func sqlNullFloat64(f float64) interface{} {
	if f == 0.0 {
		return nil
	}
	return f
}

func sqlNullTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}
