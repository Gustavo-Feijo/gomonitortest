package main

import (
	"gomonitor/internal/app"
	"gomonitor/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Error starting deps: %v", err)
	}

	log.Printf("Starting server on: %s", server.Addr)
	log.Fatal(server.Engine.Run(server.Addr))
}
