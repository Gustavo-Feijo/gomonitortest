package tracing

import (
	"context"
	"errors"
	"gomonitor/internal/config"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second),
		),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()),
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
