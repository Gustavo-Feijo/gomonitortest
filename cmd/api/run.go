package main

import (
	"context"
	"fmt"
	"gomonitor/internal/app"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/logging"
	"gomonitor/internal/observability/tracing"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context) error {
	shutdownSignalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(shutdownSignalCtx)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	logger := logging.New(cfg.Logging)

	shutdown := setupOtel(ctx, cfg.Tracing, logger)
	if shutdown != nil {
		defer shutdown()
	}

	engine, depsCleanup, err := app.New(context.Background(), cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("err", err))
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	server := &http.Server{
		Addr:    engine.Addr,
		Handler: engine.Engine,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	var serverErr error
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case serverErr = <-errCh:
		logger.Error("server stopped unexpectedly", slog.Any("err", serverErr))
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("failed to shutdown server", slog.Any("err", err))
		if serverErr == nil {
			serverErr = err
		}
	}

	if err := depsCleanup(shutdownCtx); err != nil {
		logger.Error("failed to cleanup deps", slog.Any("err", err))
		if serverErr == nil {
			serverErr = err
		}
	}

	if serverErr == nil {
		logger.Info("successfully shutdown server")
	}

	return serverErr
}

func setupOtel(ctx context.Context, cfg *config.TracingConfig, logger *slog.Logger) func() {
	// Noop otelShutdown
	otelShutdown := func(context.Context) error {
		return nil
	}

	sd, err := tracing.SetupOtel(ctx, cfg)
	if err != nil {
		logger.Warn("failed to startup otel", slog.Any("err", err))
	} else {
		otelShutdown = sd
	}

	shutdownFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otelShutdown(ctx); err != nil {
			logger.Error("failed to shutdown otel", slog.Any("err", err))
		}
	}

	return shutdownFunc
}
