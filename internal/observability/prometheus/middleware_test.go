package prometheus_test

import (
	pkgprometheus "gomonitor/internal/observability/prometheus"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func resetMetrics() {
	pkgprometheus.HTTPRequests.Reset()
	pkgprometheus.RequestDuration.Reset()
}

func TestPrometheusMiddleware_RecordsMetrics(t *testing.T) {
	resetMetrics()

	reg := prometheus.NewRegistry()
	pkgprometheus.Init(reg)

	r := gin.New()
	r.Use(pkgprometheus.PrometheusMiddleware())

	r.GET("/users", func(c *gin.Context) {
		c.Status(200)
	})

	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Must be ignored by middleware.
	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, healthReq)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	var foundCounter, foundHistogram bool

	for _, m := range metrics {
		switch m.GetName() {
		case "http_requests_total":
			foundCounter = true
			require.Len(t, m.Metric, 1)
			require.Equal(t, float64(1), m.Metric[0].Counter.GetValue())

		case "http_request_duration_seconds":
			foundHistogram = true
			require.Len(t, m.Metric, 1)
			require.Equal(t, uint64(1), m.Metric[0].Histogram.GetSampleCount())
		}
	}

	require.True(t, foundCounter)
	require.True(t, foundHistogram)
}

func TestPrometheusMiddleware_PathFallback(t *testing.T) {
	resetMetrics()

	reg := prometheus.NewRegistry()
	pkgprometheus.Init(reg)

	r := gin.New()
	r.Use(pkgprometheus.PrometheusMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	metrics, err := reg.Gather()
	require.NoError(t, err)

	found := false
	for _, m := range metrics {
		if m.GetName() == "http_requests_total" {
			found = true
			require.Len(t, m.Metric, 1)

			labels := m.Metric[0].Label
			var path string
			for _, l := range labels {
				if l.GetName() == "path" {
					path = l.GetValue()
				}
			}

			require.Equal(t, "/unknown", path)
		}
	}

	require.True(t, found)
}
