package main

import (
	"context"
	"gomonitor/internal/app"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/tracing"
	"log"
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

	otelShutdown, err := tracing.SetupOtel(ctx, cfg.Tracing)
	if err != nil {
		log.Fatalf("Error starting Otel setup: %v", err)
	}
	defer otelShutdown(ctx)

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Error starting deps: %v", err)
	}

	log.Printf("Starting server on: %s", server.Addr)
	log.Fatal(server.Engine.Run(server.Addr))
}
