// Package logger provides logging functionality for the application.
package logger

import (
	"log/slog"
	"os"
)

// NewLogger creates and returns a new structured logger instance.
func NewLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	return logger
}
