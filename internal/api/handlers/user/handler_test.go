package userhandler_test

import (
	userhandler "gomonitor/internal/api/handlers/user"
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
	handler := userhandler.NewHandler(slog.Default(), &mocks.MockUserService{}, &mocks.MockJwtManager{})

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
			name:           "create route exists",
			method:         http.MethodPost,
			path:           "/api/v1/users",
			expectedStatus: http.StatusUnauthorized,
			shouldExist:    true,
		},
		{
			name:           "get route exists",
			method:         http.MethodGet,
			path:           "/api/v1/users/1",
			expectedStatus: http.StatusUnauthorized,
			shouldExist:    true,
		},
		{
			name:           "create route only accepts POST",
			method:         http.MethodGet,
			path:           "/api/v1/users",
			expectedStatus: http.StatusMethodNotAllowed,
			shouldExist:    false,
		},
		{
			name:           "non-existent route returns 404",
			method:         http.MethodPost,
			path:           "/api/v1/nonexistent",
			expectedStatus: http.StatusNotFound,
			shouldExist:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockUserService{}
			jwtManager := &mocks.MockJwtManager{}
			h := userhandler.NewHandler(slog.Default(), mockService, jwtManager)

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
