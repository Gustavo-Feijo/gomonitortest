package authhandler_test

import (
	authhandler "gomonitor/internal/api/handlers/auth"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/mocks"
	pkgerrors "gomonitor/internal/pkg/errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
	}{
		{
			name: "service returns error",
			setupMock: func(m *mocks.MockAuthService) {
				m.On("Logout", mock.Anything).
					Return(pkgerrors.NewUnauthorizedError("Invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "successful logout",
			setupMock: func(m *mocks.MockAuthService) {
				m.On("Logout", mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockAuthService{}
			mockJwtManager := &mocks.MockJwtManager{}
			tt.setupMock(mockService)

			h := authhandler.NewHandler(slog.Default(), mockService, mockJwtManager)

			router := gin.New()
			router.HandleMethodNotAllowed = true
			router.Use(middlewares.ErrorMiddleware())
			router.POST("/logout", h.Logout)

			req := httptest.NewRequest(http.MethodPost, "/logout", http.NoBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
			mockJwtManager.AssertExpectations(t)
		})
	}
}

func TestHandler_LogoutAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
	}{
		{
			name: "service returns error",
			setupMock: func(m *mocks.MockAuthService) {
				m.On("LogoutAll", mock.Anything).
					Return(pkgerrors.NewUnauthorizedError("Invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "successful logout all",
			setupMock: func(m *mocks.MockAuthService) {
				m.On("LogoutAll", mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockAuthService{}
			mockJwtManager := &mocks.MockJwtManager{}
			tt.setupMock(mockService)

			h := authhandler.NewHandler(slog.Default(), mockService, mockJwtManager)

			router := gin.New()
			router.HandleMethodNotAllowed = true
			router.Use(middlewares.ErrorMiddleware())
			router.POST("/logoutAll", h.LogoutAll)

			req := httptest.NewRequest(http.MethodPost, "/logoutAll", http.NoBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			mockService.AssertExpectations(t)
			mockJwtManager.AssertExpectations(t)
		})
	}
}
