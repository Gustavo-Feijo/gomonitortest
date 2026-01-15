package middlewares_test

import (
	"errors"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/mocks"
	"gomonitor/internal/pkg/identity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Auth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validIdentity := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthExternal,
	}

	tests := []struct {
		name           string
		setupMock      func(*mocks.MockJwtManager)
		useQuery       bool
		authToken      string
		ginHandler     gin.HandlerFunc
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "no token",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock:      func(mjm *mocks.MockJwtManager) {},
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(mjm *mocks.MockJwtManager) {
				mjm.On("ValidateAccessToken", "invalid-token").
					Return(nil, errors.New("generic signing error"))
			},
			authToken:      "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "valid token",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(mjm *mocks.MockJwtManager) {
				mjm.On("ValidateAccessToken", "valid-token").
					Return(validIdentity, nil)
			},
			authToken:      "Bearer valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid token on queryBearer ",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(mjm *mocks.MockJwtManager) {
				mjm.On("ValidateAccessToken", "valid-token-query").
					Return(validIdentity, nil)
			},
			authToken:      "Bearer valid-token-query",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()

			jwtManagerMock := &mocks.MockJwtManager{}
			tt.setupMock(jwtManagerMock)

			// Without error middleware, auth middleware won't return correct statuses.
			r.Use(middlewares.ErrorMiddleware())

			r.Use(middlewares.AuthMiddleware(jwtManagerMock))

			r.GET("/test", tt.ginHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.useQuery {
				q := req.URL.Query()
				q.Set("token", tt.authToken)
				req.URL.RawQuery = q.Encode()
			} else {
				req.Header.Set("Authorization", tt.authToken)
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
