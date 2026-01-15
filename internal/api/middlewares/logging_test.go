package middlewares_test

import (
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/observability/logging"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Logging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	logger := slog.Default()
	r := gin.New()
	r.Use(middlewares.LoggingMiddleware(logger))

	r.GET("/test", func(c *gin.Context) {
		l := logging.FromContext(c.Request.Context())
		assert.NotNil(t, l)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
