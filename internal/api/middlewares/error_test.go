package middlewares_test

import (
	"errors"
	"gomonitor/internal/api/middlewares"
	pkgerrors "gomonitor/internal/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		ginHandler     gin.HandlerFunc
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "no errors",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "non normalized error",
			ginHandler: func(c *gin.Context) {
				_ = c.Error(errors.New("error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "normalized 500+ error",
			ginHandler: func(c *gin.Context) {
				_ = c.Error(pkgerrors.NewInternalError())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "normalized error",
			ginHandler: func(c *gin.Context) {
				_ = c.Error(pkgerrors.NewBadRequestError("bad request"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(middlewares.ErrorMiddleware())
			r.GET("/test", tt.ginHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
