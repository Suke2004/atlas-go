// Package config loads and validates all Atlas configuration from environment
// variables. Every module receives a *Config from main.go via constructor
// injection — no module reads env vars directly.
package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration. Values are populated from
// environment variables with the ATLAS_ prefix (e.g. ATLAS_PORT → Port).
// See .env.example for the full list of supported variables.
type Config struct {
	// Server
	Port int
	Env  string // "development" | "production"

	// Database
	DBPath string // path to SQLite file, e.g. ./data/db/atlas.db

	// Session
	SessionSecret string // must be exactly 32 bytes, cryptographically random
	SessionMaxAge int    // seconds; default 86400 (24h)

	// File storage
	UploadPath  string // absolute or relative path for uploaded files
	MaxUploadMB int64  // max file size in megabytes

	// AI (v2 — populated but unused until AI Workspace is built)
	AI AIConfig

	// Appearance
	ThemeDefault string // "light" | "dark" | "system"

	// Logging
	LogLevel string // "debug" | "info" | "warn" | "error"
}

// AIConfig holds provider settings for the AI backend.
// The Provider field selects which implementation is used at startup.
type AIConfig struct {
	Provider string // "openai" | "ollama" | "anthropic" | "gemini"
	APIKey   string // required for cloud providers; empty for Ollama
	BaseURL  string // e.g. http://localhost:11434 for Ollama
	Model    string // e.g. "llama3.2" or "gpt-4o"
}

// Load reads all configuration from environment variables and returns a
// populated Config. Sensible defaults are applied for every optional field.
// It calls log.Fatal if any required field (SessionSecret) is missing in
// production mode.
func Load() *Config {
	v := viper.New()

	v.SetEnvPrefix("ATLAS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Server defaults
	v.SetDefault("PORT", 8080)
	v.SetDefault("ENV", "development")

	// Database defaults
	v.SetDefault("DB_PATH", "./data/db/atlas.db")

	// Session defaults
	v.SetDefault("SESSION_SECRET", "")
	v.SetDefault("SESSION_MAX_AGE", 86400)

	// File storage defaults
	v.SetDefault("UPLOAD_PATH", "./data/uploads")
	v.SetDefault("MAX_UPLOAD_MB", 50)

	// AI defaults (all optional in v1)
	v.SetDefault("AI_PROVIDER", "ollama")
	v.SetDefault("AI_API_KEY", "")
	v.SetDefault("AI_BASE_URL", "http://localhost:11434")
	v.SetDefault("AI_MODEL", "llama3.2")

	// Appearance defaults
	v.SetDefault("THEME_DEFAULT", "system")

	// Logging defaults
	v.SetDefault("LOG_LEVEL", "info")

	cfg := &Config{
		Port:          v.GetInt("PORT"),
		Env:           v.GetString("ENV"),
		DBPath:        v.GetString("DB_PATH"),
		SessionSecret: v.GetString("SESSION_SECRET"),
		SessionMaxAge: v.GetInt("SESSION_MAX_AGE"),
		UploadPath:    v.GetString("UPLOAD_PATH"),
		MaxUploadMB:   v.GetInt64("MAX_UPLOAD_MB"),
		AI: AIConfig{
			Provider: v.GetString("AI_PROVIDER"),
			APIKey:   v.GetString("AI_API_KEY"),
			BaseURL:  v.GetString("AI_BASE_URL"),
			Model:    v.GetString("AI_MODEL"),
		},
		ThemeDefault: v.GetString("THEME_DEFAULT"),
		LogLevel:     v.GetString("LOG_LEVEL"),
	}

	// Enforce required fields in production
	if cfg.Env == "production" && cfg.SessionSecret == "" {
		log.Fatal("ATLAS_SESSION_SECRET is required in production mode. " +
			"Generate one with: openssl rand -hex 32")
	}

	return cfg
}

// IsDevelopment returns true when running in development mode.
// Used to enable debug logging, pretty-printed logs, etc.
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true when running in production mode.
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
