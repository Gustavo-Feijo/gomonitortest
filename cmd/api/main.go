package main

import (
	"context"
	"gomonitor/internal/app"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/logging"
	"gomonitor/internal/observability/tracing"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := logging.New(cfg.Logging)

	otelShutdown, err := tracing.SetupOtel(ctx, cfg.Tracing)
	if err != nil {
		logger.Error("failed to startup otel", slog.Any("err", err))
		return
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			logger.Error("failed to shutdown otel", slog.Any("err", err))
		}
	}()

	server, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("err", err))
		return
	}

	logger.Info("starting server", slog.String("addr", server.Addr))
	if err := server.Engine.Run(server.Addr); err != nil {
		logger.Error("server stopped", slog.Any("err", err))
	}
}
