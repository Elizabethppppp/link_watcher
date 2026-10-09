package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level  string
	Format string
}

func Init(cfg Config) {
	level := parseLevel(cfg.Level)

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}

	var handler slog.Handler
	if strings.EqualFold(cfg.Format, "json") {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func parseLevel(l string) slog.Level {
	switch strings.ToLower(l) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Info(msg string, v ...any) {
	slog.Info(msg, v...)
}

func Warn(msg string, v ...any) {
	slog.Warn(msg, v...)
}

func Error(msg string, v ...any) {
	slog.Error(msg, v...)
}

func Debug(msg string, v ...any) {
	slog.Debug(msg, v...)
}

func Fatal(msg string, v ...any) {
	slog.Error(msg, v...)
	os.Exit(1)
}
