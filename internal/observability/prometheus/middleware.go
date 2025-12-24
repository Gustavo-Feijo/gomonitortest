package prometheus

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()

		if path == "/health" || path == "/metrics" {
			return
		}

		if path == "" {
			path = c.Request.URL.Path
		}

		HTTPRequests.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()

		RequestDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Observe(time.Since(start).Seconds())
	}
}
