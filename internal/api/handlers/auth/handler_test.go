package authhandler_test

import (
	authhandler "gomonitor/internal/api/handlers/auth"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/mocks"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NewHandler(t *testing.T) {
	handler := authhandler.NewHandler(slog.Default(), &mocks.MockAuthService{})

	assert.NotNil(t, handler)
}

func TestHandler_RegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		shouldExist    bool
	}{
		{
			name:           "login route exists",
			method:         http.MethodPost,
			path:           "/api/v1/auth/login",
			expectedStatus: http.StatusBadRequest,
			shouldExist:    true,
		},
		{
			name:           "refresh route exists",
			method:         http.MethodPost,
			path:           "/api/v1/auth/refresh",
			expectedStatus: http.StatusBadRequest,
			shouldExist:    true,
		},
		{
			name:           "login only accepts POST",
			method:         http.MethodGet,
			path:           "/api/v1/auth/login",
			expectedStatus: http.StatusMethodNotAllowed,
			shouldExist:    false,
		},
		{
			name:           "non-existent route returns 404",
			method:         http.MethodPost,
			path:           "/api/v1/auth/logout",
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockAuthService{}
			h := authhandler.NewHandler(slog.Default(), mockService)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.HandleMethodNotAllowed = true
			router.Use(middlewares.ErrorMiddleware())

			apiV1 := router.Group("/api/v1")

			h.RegisterRoutes(apiV1)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
