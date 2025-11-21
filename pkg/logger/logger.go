package logger

import (
	"log/slog"
	"os"
)

var Global *slog.Logger

// Init инициализирует глобальный логгер
func Init() {
	Global = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	Global.Info("logger initialized")
}

// Shortcut функции для удобства
func Info(msg string, args ...interface{}) {
	Global.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Global.Error(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Global.Warn(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Global.Debug(msg, args...)
}
