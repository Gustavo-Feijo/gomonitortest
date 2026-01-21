package middlewares

import (
	pkgprometheus "gomonitor/internal/observability/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware() gin.HandlerFunc {
	pkgprometheus.Init(prometheus.DefaultRegisterer)
	return pkgprometheus.PrometheusMiddleware()
}
