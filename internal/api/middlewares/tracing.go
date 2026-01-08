package middlewares

import (
	"gomonitor/internal/config"
	"gomonitor/internal/observability/tracing"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func TracingMiddleware(cfg *config.Config) gin.HandlerFunc {
	return otelgin.Middleware(cfg.Tracing.ServiceName,
		otelgin.WithFilter(func(r *http.Request) bool {
			return !slices.Contains(tracing.IgnoredRoutes, r.URL.Path)
		}),
	)
}
