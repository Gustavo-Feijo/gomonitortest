package userhandler_test

import (
	"encoding/json"
	userdto "gomonitor/internal/api/dto/user"
	userhandler "gomonitor/internal/api/handlers/user"
	"gomonitor/internal/api/middlewares"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/mocks"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	defaultUser := userdto.GetUserRequest{
		ID: 1,
	}

	now := time.Now()

	userReturn := &user.User{
		ID:        1,
		Name:      "test",
		Email:     "test@example.com",
		UserName:  "test",
		Password:  "generated-hash",
		Role:      identity.RoleUser,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name           string
		route          string
		setupMock      func(*mocks.MockUserService)
		expectedStatus int
		validateResp   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "invalid path",
			route:          "/invalid",
			setupMock:      func(m *mocks.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "Invalid ID paramete")
			},
		},
		{
			name:  "service returns error",
			route: "/1",
			setupMock: func(m *mocks.MockUserService) {
				m.On("GetUser", mock.Anything, defaultUser.ToDomainInput()).
					Return(nil, pkgerrors.NewUnauthorizedError("unauthenticated"))
			},
			expectedStatus: http.StatusUnauthorized,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				assert.Contains(t, rec.Body.String(), "unauthenticated")
			},
		},
		{
			name:  "successful retrieve",
			route: "/1",
			setupMock: func(m *mocks.MockUserService) {
				m.On("GetUser", mock.Anything, defaultUser.ToDomainInput()).
					Return(userReturn, nil)
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp userdto.GetUserResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				assert.Equal(t, userReturn.ID, resp.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockUserService{}
			mockJwt := &mocks.MockJwtManager{}
			tt.setupMock(mockService)

			h := userhandler.NewHandler(slog.Default(), mockService, mockJwt)

			router := gin.New()
			router.HandleMethodNotAllowed = true
			router.Use(middlewares.ErrorMiddleware())
			router.POST("/:id", h.GetByID)

			req := httptest.NewRequest(http.MethodPost, tt.route, nil)
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
