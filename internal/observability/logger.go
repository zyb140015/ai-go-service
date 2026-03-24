package observability

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a structured logger with the configured level.
func NewLogger(level string) (*slog.Logger, error) {
	var slogLevel slog.Level

	switch strings.ToUpper(level) {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		return nil, fmt.Errorf("unsupported log level %q", level)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})

	return slog.New(handler), nil
}
