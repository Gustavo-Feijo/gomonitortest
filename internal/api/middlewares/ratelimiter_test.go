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
	"github.com/stretchr/testify/mock"
)

type ratelimitTestMocks struct {
	rateLimiterMock *mocks.MockRateLimiter
}

func TestMiddleware_RateLimitIp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*ratelimitTestMocks)
		ginHandler     gin.HandlerFunc
		expectedStatus int
	}{
		{
			name: "ratelimiter error",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(false, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "not allowed",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(false, nil)
			},
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name: "allowed",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(true, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rateLimiterMock := &mocks.MockRateLimiter{}

			ratelimitTestMocks := &ratelimitTestMocks{
				rateLimiterMock: rateLimiterMock,
			}

			tt.setupMock(ratelimitTestMocks)

			// Without error middleware, ratelimiter middleware won't return correct statuses.
			r.Use(middlewares.ErrorMiddleware())

			r.Use(middlewares.IPRateLimiterMiddleware(rateLimiterMock))

			r.GET("/test", tt.ginHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			rateLimiterMock.AssertExpectations(t)
		})
	}
}

func TestMiddleware_RateLimitUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validIdentity := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthExternal,
	}

	tests := []struct {
		name           string
		setupMock      func(*ratelimitTestMocks)
		principal      *identity.Principal
		ginHandler     gin.HandlerFunc
		expectedStatus int
	}{
		{
			name: "no principal",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			setupMock:      func(rlm *ratelimitTestMocks) {},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ratelimiter error",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			principal: validIdentity,
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(false, errors.New("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "not allowed",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			principal: validIdentity,
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(false, nil)
			},
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name: "allowed",
			ginHandler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			principal: validIdentity,
			setupMock: func(rlm *ratelimitTestMocks) {
				rlm.rateLimiterMock.
					On("Allow", mock.Anything, mock.Anything).
					Return(true, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			rateLimiterMock := &mocks.MockRateLimiter{}

			ratelimitTestMocks := &ratelimitTestMocks{
				rateLimiterMock: rateLimiterMock,
			}

			tt.setupMock(ratelimitTestMocks)

			// Without error middleware, ratelimiter middleware won't return correct statuses.
			r.Use(middlewares.ErrorMiddleware())

			r.Use(middlewares.UserRateLimiterMiddleware(rateLimiterMock))

			r.GET("/test", tt.ginHandler)

			ctx := t.Context()
			if tt.principal != nil {
				ctx = identity.WithPrincipal(ctx, tt.principal)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			rateLimiterMock.AssertExpectations(t)
		})

	}
}
