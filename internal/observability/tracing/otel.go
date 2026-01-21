package tracing

import (
	"context"
	"errors"
	"gomonitor/internal/config"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var IgnoredRoutes = []string{
	"/health",
	"/metrics",
}

// SetupOtel starts the Otel.
func SetupOtel(ctx context.Context, cfg *config.TracingConfig) (func(context.Context) error, error) {
	traceExporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(cfg.Address),
			otlptracegrpc.WithInsecure(),
		),
	)
	if err != nil {
		if otelShutdownErr := traceExporter.Shutdown(ctx); otelShutdownErr != nil {
			err = errors.Join(err, otelShutdownErr)
		}
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		if otelShutdownErr := traceExporter.Shutdown(ctx); otelShutdownErr != nil {
			err = errors.Join(err, otelShutdownErr)
		}
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBatchTimeout(time.Second),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	otel.SetTracerProvider(tracerProvider)

	return traceExporter.Shutdown, err
}

// Tracer returns the global tracer
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// RecordError records an error in the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
}

// SetAttributes adds attributes to the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, tracerName, spanName string) (context.Context, trace.Span) {
	tracer := Tracer(tracerName)
	return tracer.Start(ctx, spanName)
}

// TrackDuration tracks the duration of an operation
func TrackDuration(ctx context.Context, operation string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		SetAttributes(ctx,
			attribute.String("operation", operation),
			attribute.Float64("duration_ms", float64(duration.Milliseconds())),
		)
	}
}
