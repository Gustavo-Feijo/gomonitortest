package authhandler_test

import (
	"bytes"
	"encoding/json"
	authdto "gomonitor/internal/api/dto/auth"
	authhandler "gomonitor/internal/api/handlers/auth"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/mocks"
	pkgerrors "gomonitor/internal/pkg/errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    any
		setupMock      func(*mocks.MockAuthService)
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "invalid JSON payload",
			requestBody:    "invalidjson",
			setupMock:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Invalid JSON payload")
			},
		},
		{
			name:           "missing required fields",
			requestBody:    authdto.RefreshRequest{},
			setupMock:      func(m *mocks.MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			validateResp:   func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name: "service returns error",
			requestBody: authdto.RefreshRequest{
				RefreshToken: "refresh-token",
			},
			setupMock: func(m *mocks.MockAuthService) {
				m.On("Refresh", mock.Anything, mock.Anything).
					Return(nil, pkgerrors.NewUnauthorizedError("Invalid credentials", nil))
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Invalid credentials")
			},
		},
		{
			name: "successful refresh",
			requestBody: authdto.RefreshRequest{
				RefreshToken: "refresh-token",
			},
			setupMock: func(m *mocks.MockAuthService) {
				m.On("Refresh", mock.Anything, mock.MatchedBy(func(input auth.RefreshInput) bool {
					return input.RefreshToken == "refresh-token"
				})).Return(&auth.RefreshOutput{
					AccessToken: "access-token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp authdto.RefreshResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, "access-token", resp.AccessToken)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockAuthService{}
			tt.setupMock(mockService)

			h := authhandler.NewHandler(slog.Default(), mockService)

			router := gin.New()
			router.HandleMethodNotAllowed = true
			router.Use(middlewares.ErrorMiddleware())
			router.POST("/refresh", h.Refresh)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.validateResp != nil {
				tt.validateResp(t, rec)
			}

			mockService.AssertExpectations(t)
		})
	}
}
