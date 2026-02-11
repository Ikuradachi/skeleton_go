package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init initializes the default logger
func Init() {
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Get returns the default logger instance
func Get() *slog.Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}

// NewLogger creates a new logger with custom options
func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

// Info logs an info message
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}
