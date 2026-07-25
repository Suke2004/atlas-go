// Package logger wraps go.uber.org/zap and provides a single constructor that
// builds the correct encoder (console for dev, JSON for production) based on
// the environment setting from Config.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a configured *zap.Logger.
//
// In development: human-readable console output, DEBUG level.
// In production:  structured JSON output, INFO level (or as configured).
//
// The caller is responsible for calling log.Sync() on shutdown.
func New(env, level string) *zap.Logger {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Apply configured log level
	if lvl, err := zap.ParseAtomicLevel(level); err == nil {
		cfg.Level = lvl
	}

	log, err := cfg.Build()
	if err != nil {
		// Fall back to a no-op logger rather than crashing at startup.
		return zap.NewNop()
	}

	return log
}
