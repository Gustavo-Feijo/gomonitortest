package middlewares

import (
	"gomonitor/internal/infra/deps"
	"gomonitor/internal/observability/logging"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(deps *deps.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger := logging.WithTrace(c.Request.Context(), deps.Logger)
		ctx := logging.WithContext(c.Request.Context(), ctxLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
