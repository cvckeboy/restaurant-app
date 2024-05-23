package utils

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

// NewLogger initializes and returns a new logger instance.
func NewLogger(config *Config) *Logger {
	var logLevel slog.Level
	var addSource bool
	switch config.Logging.Level {
	case "debug":
		logLevel = slog.LevelDebug
		addSource = true
	case "info":
		logLevel = slog.LevelInfo
		addSource = true
	case "warn":
		logLevel = slog.LevelWarn
		addSource = false
	case "error":
		logLevel = slog.LevelError
		addSource = false
	default:
		logLevel = slog.LevelInfo
		addSource = true
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel, AddSource: addSource})
	return &Logger{slog.New(handler)}
}
