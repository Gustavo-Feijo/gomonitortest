package middlewares

import (
	"gomonitor/internal/observability/logging"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxLogger := logging.WithTrace(c.Request.Context(), logger)
		ctx := logging.WithContext(c.Request.Context(), ctxLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
