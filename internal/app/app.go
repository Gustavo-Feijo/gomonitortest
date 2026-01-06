package app

import (
	"context"
	"fmt"
	"gomonitor/internal/config"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/domain/user"
	databaseinfra "gomonitor/internal/infra/database"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/observability/logging"
	"gomonitor/internal/observability/prometheus"
	"gomonitor/internal/observability/tracing"
	"log/slog"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	engine.Use(gin.Recovery())

	// Setup prometheus middleware.
	prometheus.Init()
	engine.Use(prometheus.PrometheusMiddleware())

	// Setup tracing middleware.
	tracingMiddleware := otelgin.Middleware(cfg.Tracing.ServiceName,
		otelgin.WithFilter(func(r *http.Request) bool {
			return !slices.Contains(tracing.IgnoredRoutes, r.URL.Path)
		}),
	)
	engine.Use(tracingMiddleware)

	// Setup logging to store trace and span ids..
	engine.Use(
		func(c *gin.Context) {
			ctxLogger := logging.WithTrace(c.Request.Context(), deps.Logger)
			c.Set("logger", ctxLogger)
			c.Next()
		},
	)

	userHandler := user.NewHandler(deps)
	authHandler := auth.NewHandler(deps, cfg.Auth)

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
