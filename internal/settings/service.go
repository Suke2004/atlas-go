package settings

import (
	"context"
	"errors"

	"github.com/Suke2004/atlas-go/internal/db"
)

// Service manages user settings stored in the key-value settings table.
type Service interface {
	Get(ctx context.Context, userID int64, key string) string
	Set(ctx context.Context, userID int64, key, value string) error
	GetAll(ctx context.Context, userID int64) map[string]string
}

type service struct {
	database *db.DB
}

// NewService creates a settings service backed by the SQLite settings table.
func NewService(database *db.DB) Service {
	return &service{database: database}
}

// Get returns the value for the given key, or "" if not found.
func (s *service) Get(ctx context.Context, userID int64, key string) string {
	val, err := s.database.Queries.GetSetting(ctx, db.GetSettingParams{
		UserID: userID,
		Key:    key,
	})
	if err != nil {
		return ""
	}
	return val
}

// Set upserts a setting value.
func (s *service) Set(ctx context.Context, userID int64, key, value string) error {
	if key == "" {
		return errors.New("settings: key cannot be empty")
	}
	return s.database.Queries.SetSetting(ctx, db.SetSettingParams{
		UserID: userID,
		Key:    key,
		Value:  value,
	})
}

// GetAll returns all settings for the user as a map.
func (s *service) GetAll(ctx context.Context, userID int64) map[string]string {
	rows, err := s.database.Queries.GetAllSettings(ctx, userID)
	out := map[string]string{
		// Sensible defaults so the UI never shows empty fields
		"theme":            "dark",
		"ai_provider":      "ollama",
		"ai_ollama_url":    "http://localhost:11434",
		"ai_ollama_model":  "llama3.2",
		"ai_gemini_model":  "gemini-2.0-flash",
		"ai_openai_model":  "gpt-4o-mini",
	}
	if err != nil {
		return out
	}
	for _, r := range rows {
		out[r.Key] = r.Value
	}
	return out
}
