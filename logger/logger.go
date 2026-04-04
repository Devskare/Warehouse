package logger

import (
	"log/slog"
	"os"
)

func Initlogger(Level string, Env bool) *slog.Logger {
	var level slog.Level

	switch Level {
	case "debug":
		level = slog.LevelDebug
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if Env == true {
		handler := slog.NewJSONHandler(os.Stdout, opts)
		return slog.New(handler)
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
