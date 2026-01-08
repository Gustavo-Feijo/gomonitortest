package middlewares

import (
	"gomonitor/internal/observability/prometheus"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	prometheus.Init()
	return prometheus.PrometheusMiddleware()
}
