package middlewares_test

import (
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Tracing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	tracingMiddleware := &config.Config{
		Tracing: &config.TracingConfig{
			ServiceName: "test",
		},
	}

	assert.NotNil(t, middlewares.TracingMiddleware(tracingMiddleware))
	r.Use(middlewares.TracingMiddleware(tracingMiddleware))

	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest("GET", "/ping", nil))
	assert.Equal(t, http.StatusOK, w.Code)
}
