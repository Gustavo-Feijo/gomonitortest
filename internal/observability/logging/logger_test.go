package logging

import (
	"bytes"
	"context"
	"gomonitor/internal/config"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestNewLogger(t *testing.T) {
	logger := New(&config.LoggingConfig{Level: "debug"})
	require.NotNil(t, logger)
}

func TestNewLoggerDefaultLevel(t *testing.T) {
	logger := New(&config.LoggingConfig{Level: "xxxx"})
	require.NotNil(t, logger)
}

func TestWithContext(t *testing.T) {
	logger := New(&config.LoggingConfig{Level: "error"})
	require.NotNil(t, logger)

	ctx := context.Background()
	newCtx := WithContext(ctx, logger)
	assert.NotNil(t, newCtx.Value(loggerKey))
}

func TestFromContext(t *testing.T) {
	logger := New(&config.LoggingConfig{Level: "warn"})
	require.NotNil(t, logger)

	ctx := context.Background()
	newCtx := WithContext(ctx, logger)
	assert.NotNil(t, newCtx.Value(loggerKey))

	loggerFromCtx := FromContext(newCtx)
	require.NotNil(t, loggerFromCtx)
}

func TestFromContextFallback(t *testing.T) {
	ctx := context.Background()
	loggerFromCtx := FromContext(ctx)
	require.NotNil(t, loggerFromCtx)
	loggerFromCtx.Info("hello")
}

func TestWithTrace_NoSpan(t *testing.T) {
	var buf bytes.Buffer
	logger := testLogger(&buf)

	ctx := context.Background()

	newLogger := WithTrace(ctx, logger)

	require.Same(t, logger, newLogger)

	newLogger.Info("hello")

	output := buf.String()
	require.NotContains(t, output, "trace_id")
	require.NotContains(t, output, "span_id")
}

func TestWithTrace_WithSpan(t *testing.T) {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	defer func() { _ = tp.Shutdown(context.Background()) }()

	var buf bytes.Buffer
	logger := testLogger(&buf)

	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "span")
	defer span.End()

	newLogger := WithTrace(ctx, logger)

	require.NotSame(t, logger, newLogger)

	newLogger.Info("hello")

	output := buf.String()
	assert.Contains(t, output, "trace_id=")
	assert.Contains(t, output, "span_id=")

	sc := span.SpanContext()
	assert.Contains(t, output, sc.TraceID().String())
	assert.Contains(t, output, sc.SpanID().String())
}

func testLogger(buf *bytes.Buffer) *slog.Logger {
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}
