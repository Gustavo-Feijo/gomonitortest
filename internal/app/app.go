package app

import (
	"gomonitor/internal/config"
	databaseinfra "gomonitor/internal/infra/database"
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/observability/prometheus"
	"gomonitor/internal/user"

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
func New(cfg *config.Config) (*App, error) {
	deps, err := deps.NewDeps(cfg)
	if err != nil {
		return nil, err
	}

	if err := databaseinfra.RunMigrations(cfg.Database, deps.DB); err != nil {
		return nil, err
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	prometheus.Init()
	engine.Use(prometheus.PrometheusMiddleware())

	engine.Use(otelgin.Middleware(cfg.Tracing.ServiceName))

	userHander := user.NewHandler(deps)
	registerRoutes(engine, userHander)

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
