package logging

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func NewLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.SourceKey:
				if s, ok := a.Value.Any().(*slog.Source); ok {
					a.Value = slog.StringValue(fmt.Sprintf("%s:%d", filepath.Base(s.File), s.Line))
				}
			case slog.TimeKey:
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05"))
				}
			}
			return a
		},
	})

	return slog.New(handler)
}
