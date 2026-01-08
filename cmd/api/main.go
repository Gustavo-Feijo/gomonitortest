package main

import (
	"context"
	"gomonitor/internal/app"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/logging"
	"gomonitor/internal/observability/tracing"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	shutdownSignalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithCancel(shutdownSignalCtx)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	logger := logging.New(cfg.Logging)

	otelShutdown, err := tracing.SetupOtel(context.Background(), cfg.Tracing)
	if err != nil {
		logger.Error("failed to startup otel", slog.Any("err", err))
		return
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otelShutdown(ctx); err != nil {
			logger.Error("failed to shutdown otel", slog.Any("err", err))
		}

	}()

	engine, err := app.New(context.Background(), cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("err", err))
		return
	}

	server := &http.Server{
		Addr:    engine.Addr,
		Handler: engine.Engine,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped", slog.Any("err", err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case err := <-errCh:
		logger.Error("server stopped unexpectedly", slog.Any("err", err))
		cancel()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server", slog.Any("err", err))
	}
	logger.Info("successfully shutdown server")
}
