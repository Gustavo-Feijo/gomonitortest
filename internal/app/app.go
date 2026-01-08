package app

import (
	"context"
	"fmt"
	authhandler "gomonitor/internal/api/handlers/auth"
	userhandler "gomonitor/internal/api/handlers/user"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	"gomonitor/internal/infra/deps"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App consists of the engine and address where it will listen.
type App struct {
	Engine *gin.Engine
	Addr   string
}

// RouteRegister is a interface to be implemented by the handlers to define routes.
type RouteRegister interface {
	RegisterRoutes(r *gin.RouterGroup)
}

// New returns a new app.
func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	deps, err := deps.NewDeps(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("error initializing dependencies: %v", err)
	}

	if err := databaseinfra.RunMigrations(ctx, cfg.Database, deps.DB); err != nil {
		return nil, fmt.Errorf("error running migrations: %v", err)
	}

	if err = bootstrapApp(ctx, cfg, deps); err != nil {
		logger.Error("error bootstrapping application", slog.Any("err", err))
	}

	engine := gin.New()

	// Add middlewares.
	engine.Use(gin.Recovery())
	engine.Use(middlewares.TracingMiddleware(cfg))
	engine.Use(middlewares.LoggingMiddleware(deps))
	engine.Use(middlewares.ErrorMiddleware())
	engine.Use(middlewares.PrometheusMiddleware())

	userHandler := userhandler.NewHandler(deps)
	authHandler := authhandler.NewHandler(deps, cfg.Auth)

	registerRoutes(engine, userHandler, authHandler)

	return &App{
		Engine: engine,
		Addr:   cfg.HTTP.Address,
	}, nil
}

// registerRoutes defines the base API path and calls the  handlers.
func registerRoutes(r *gin.Engine, handlers ...RouteRegister) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api")

	for _, h := range handlers {
		h.RegisterRoutes(api)
	}
}
