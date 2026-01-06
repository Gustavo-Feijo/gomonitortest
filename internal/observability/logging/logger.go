package logging

import (
	"context"
	"gomonitor/internal/config"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// New creates a new structured logger.
func New(cfg *config.LoggingConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Level),
		AddSource: cfg.Level == "debug",
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}

// WithTrace adds a trace/span id to the logger.
func WithTrace(ctx context.Context, logger *slog.Logger) *slog.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return logger
	}

	return logger.With(
		slog.String("trace_id", span.SpanContext().TraceID().String()),
		slog.String("span_id", span.SpanContext().SpanID().String()),
	)
}

// FromContext is a helper to extract the logger from the Gin context.
func FromContext(c *gin.Context) *slog.Logger {
	if l, exists := c.Get("logger"); exists {
		if logger, ok := l.(*slog.Logger); ok && logger != nil {
			return logger
		}
	}
	return slog.Default()
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
