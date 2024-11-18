package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Level
	switch env {
	case "local":
		level = slog.LevelDebug
	case "dev":
		level = slog.LevelInfo
	default:
		level = slog.LevelWarn
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	logger.Info("Logger initialized", "environment", env, "level", level.String())
	return logger
}

