package tracing_test

import (
	"context"
	"errors"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/tracing"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestSetupOtel_Success(t *testing.T) {
	ctx := context.Background()

	cfg := &config.TracingConfig{
		Address:     "localhost:4317",
		ServiceName: "test-service",
	}

	shutdown, err := tracing.SetupOtel(ctx, cfg)

	require.NoError(t, err)
	require.NotNil(t, shutdown)

	tracer := otel.Tracer("test")
	require.NotNil(t, tracer)

	err = shutdown(ctx)
	require.NoError(t, err)
}

func setupTestTracer(t *testing.T) (*tracetest.InMemoryExporter, func()) {
	t.Helper()

	exporter := tracetest.NewInMemoryExporter()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.Empty()),
	)

	otel.SetTracerProvider(tp)

	return exporter, func() {
		_ = tp.Shutdown(context.Background())
	}
}

func TestTracer(t *testing.T) {
	_, shutdown := setupTestTracer(t)
	defer shutdown()

	tr := tracing.Tracer("test-tracer")
	require.NotNil(t, tr)
}

func TestStartSpan(t *testing.T) {
	exporter, shutdown := setupTestTracer(t)
	defer shutdown()

	ctx := context.Background()

	_, span := tracing.StartSpan(ctx, "test-tracer", "test-span")
	require.NotNil(t, span)

	span.End()

	spans := exporter.GetSpans()
	require.Len(t, spans, 1)
	require.Equal(t, "test-span", spans[0].Name)
}

func TestSetAttributes(t *testing.T) {
	exporter, shutdown := setupTestTracer(t)
	defer shutdown()

	ctx, span := tracing.StartSpan(context.Background(), "test", "attr-span")

	tracing.SetAttributes(ctx,
		attribute.String("key", "value"),
		attribute.Int("answer", 42),
	)

	span.End()

	spans := exporter.GetSpans()
	require.Len(t, spans, 1)

	attrs := spans[0].Attributes

	hasAttr := func(key string) bool {
		for _, a := range attrs {
			if string(a.Key) == key {
				return true
			}
		}
		return false
	}

	require.True(t, hasAttr("key"))
	require.True(t, hasAttr("answer"))
}

func TestRecordError(t *testing.T) {
	exporter, shutdown := setupTestTracer(t)
	defer shutdown()

	ctx, span := tracing.StartSpan(context.Background(), "test", "error-span")

	tracing.RecordError(ctx, errors.New("boom"))

	span.End()

	spans := exporter.GetSpans()
	require.Len(t, spans, 1)
	require.NotEmpty(t, spans[0].Events)

	found := false
	for _, e := range spans[0].Events {
		if e.Name == "exception" {
			found = true
			break
		}
	}

	require.True(t, found, "expected exception event")
}

func TestTrackDuration(t *testing.T) {
	exporter, shutdown := setupTestTracer(t)
	defer shutdown()

	ctx, span := tracing.StartSpan(context.Background(), "test", "duration-span")

	done := tracing.TrackDuration(ctx, "test-operation")
	time.Sleep(5 * time.Millisecond)
	done()

	span.End()

	spans := exporter.GetSpans()
	require.Len(t, spans, 1)

	found := false
	for _, a := range spans[0].Attributes {
		if string(a.Key) == "duration_ms" {
			found = true
			break
		}
	}

	require.True(t, found, "duration_ms attribute not found")
}
