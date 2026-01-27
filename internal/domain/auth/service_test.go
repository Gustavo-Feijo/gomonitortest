package auth_test

import (
	"context"
	"errors"
	"gomonitor/internal/config"
	"gomonitor/internal/domain/auth"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/domain/user/testdata"
	"gomonitor/internal/mocks"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/jwt"
	"gomonitor/internal/testutil"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginMocks struct {
	userRepo         *mocks.MockUserRepository
	refreshTokenRepo *mocks.MockRefreshTokenRepository
	hasher           *mocks.MockPasswordHasher
	jwtManager       *mocks.MockJwtManager
}

type logoutMocks struct {
	refreshTokenRepo *mocks.MockRefreshTokenRepository
}

type refreshMocks struct {
	userRepo         *mocks.MockUserRepository
	refreshTokenRepo *mocks.MockRefreshTokenRepository
	jwtManager       *mocks.MockJwtManager
}

func TestNewService(t *testing.T) {
	t.Parallel()
	deps := &auth.ServiceDeps{
		UserRepo:     user.NewUserRepository(&gorm.DB{}),
		Logger:       &slog.Logger{},
		TokenManager: jwt.NewTokenManager(&config.AuthConfig{}),
		AuthConfig:   &config.AuthConfig{},
	}

	service := auth.NewService(deps)
	assert.NotNil(t, service)
}

func TestService_Login(t *testing.T) {
	t.Parallel()

	fakeHash := testdata.TestPasswordHash

	defaultInput := auth.LoginInput{
		Email:    "test@test.com",
		Password: "password123",
	}

	defaultJti := uuid.New()
	defaultUserReturn := &user.User{
		ID:       1,
		Name:     "test",
		UserName: "test",
		Email:    "test@test.com",
		Password: fakeHash,
		Role:     identity.RoleUser,
	}

	fakeRefreshToken := "fakeRefresh"
	fakeAccessToken := "fakeAccess"

	fakeRefreshTokenResult := &jwt.RefreshTokenResult{
		Token: fakeRefreshToken,
		Meta: jwt.TokenMetadata{
			JTI: defaultJti,
		},
	}

	fakeAccessTokenResult := &jwt.AccessTokenResult{
		Token: fakeAccessToken,
		Meta: jwt.TokenMetadata{
			JTI: defaultJti,
		},
	}

	tests := []struct {
		name       string
		input      auth.LoginInput
		setupMocks func(m *loginMocks)
		expected   *auth.LoginOutput
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "non existing user",
			input: auth.LoginInput{
				Email:    "nonexistenttest@test.com",
				Password: "password123",
			},
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "nonexistenttest@test.com").
					Return(testutil.Err[*user.User](gorm.ErrRecordNotFound))

				m.hasher.
					On("VerifyPassword", fakeHash, "password123").
					Return(nil)
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "db error",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Err[*user.User](gorm.ErrInvalidDB))

				m.hasher.
					On("VerifyPassword", fakeHash, "password123").
					Return(nil)
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "wrong password",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				m.hasher.
					On("VerifyPassword", fakeHash, "password123").
					Return(bcrypt.ErrMismatchedHashAndPassword)
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "refresh token error",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				m.hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				m.jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(nil, errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "refresh token store error",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				m.hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				m.jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeRefreshTokenResult, nil)

				m.refreshTokenRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(errors.New("error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "access token error",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				m.hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				m.jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeRefreshTokenResult, nil)

				m.refreshTokenRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(nil)

				m.jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role, mock.Anything).
					Return(nil, errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "success",
			input: defaultInput,
			setupMocks: func(m *loginMocks) {
				m.userRepo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				m.hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				m.jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeRefreshTokenResult, nil)

				m.refreshTokenRepo.
					On("Create", mock.Anything, mock.Anything).
					Return(nil)

				m.jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role, mock.Anything).
					Return(fakeAccessTokenResult, nil)
			},
			expected: &auth.LoginOutput{
				RefreshToken: fakeRefreshToken,
				AccessToken:  fakeAccessToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := &mocks.MockUserRepository{}
			refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
			hasher := &mocks.MockPasswordHasher{}
			jwtManager := &mocks.MockJwtManager{}

			loginMocks := &loginMocks{
				userRepo:         userRepo,
				refreshTokenRepo: refreshTokenRepo,
				hasher:           hasher,
				jwtManager:       jwtManager,
			}

			tt.setupMocks(loginMocks)

			svcDeps := &auth.ServiceDeps{
				AuthConfig: &config.AuthConfig{
					FakeHash: fakeHash,
				},
				Hasher:           hasher,
				UserRepo:         userRepo,
				Logger:           slog.Default(),
				TokenManager:     jwtManager,
				RefreshTokenRepo: refreshTokenRepo,
			}
			service := auth.NewService(svcDeps)

			result, err := service.Login(t.Context(), tt.input)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			userRepo.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			jwtManager.AssertExpectations(t)
		})
	}
}

func TestService_Logout(t *testing.T) {
	t.Parallel()

	defaultJti := uuid.New()

	tests := []struct {
		name       string
		setupCtx   func(ctx context.Context) context.Context
		setupMocks func(m *logoutMocks)
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "no principal",
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "missing refresh jti",
			setupCtx: func(ctx context.Context) context.Context {
				return identity.WithPrincipal(ctx, &identity.Principal{
					UserID: 1,
					Role:   identity.RoleAdmin,
					Source: identity.AuthExternal,
				})
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "db revoking error",
			setupMocks: func(m *logoutMocks) {
				m.refreshTokenRepo.
					On("RevokeByJTI", mock.Anything, defaultJti).
					Return(errors.New("error"))
			},
			setupCtx: func(ctx context.Context) context.Context {
				return identity.WithPrincipal(ctx, &identity.Principal{
					UserID:     1,
					Role:       identity.RoleAdmin,
					Source:     identity.AuthExternal,
					RefreshJTI: &defaultJti,
				})
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "success",
			setupMocks: func(m *logoutMocks) {
				m.refreshTokenRepo.
					On("RevokeByJTI", mock.Anything, defaultJti).
					Return(nil)
			},
			setupCtx: func(ctx context.Context) context.Context {
				return identity.WithPrincipal(ctx, &identity.Principal{
					UserID:     1,
					Role:       identity.RoleAdmin,
					Source:     identity.AuthExternal,
					RefreshJTI: &defaultJti,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshTokenRepo := &mocks.MockRefreshTokenRepository{}

			logoutMocks := &logoutMocks{
				refreshTokenRepo: refreshTokenRepo,
			}

			if tt.setupMocks != nil {
				tt.setupMocks(logoutMocks)
			}

			svcDeps := &auth.ServiceDeps{
				Logger:           slog.Default(),
				RefreshTokenRepo: refreshTokenRepo,
			}
			service := auth.NewService(svcDeps)

			ctx := t.Context()
			if tt.setupCtx != nil {
				ctx = tt.setupCtx(ctx)
			}

			err := service.Logout(ctx)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
			} else {
				assert.NoError(t, err)
			}

			refreshTokenRepo.AssertExpectations(t)
		})
	}
}

func TestService_LogoutAll(t *testing.T) {
	t.Parallel()

	userId := uint(1)
	tests := []struct {
		name       string
		setupCtx   func(ctx context.Context) context.Context
		setupMocks func(m *logoutMocks)
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "no principal",
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "db revoking error",
			setupMocks: func(m *logoutMocks) {
				m.refreshTokenRepo.
					On("RevokeByUserID", mock.Anything, userId).
					Return(errors.New("error"))
			},
			setupCtx: func(ctx context.Context) context.Context {
				return identity.WithPrincipal(ctx, &identity.Principal{
					UserID: userId,
					Role:   identity.RoleAdmin,
					Source: identity.AuthExternal,
				})
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "success",
			setupMocks: func(m *logoutMocks) {
				m.refreshTokenRepo.
					On("RevokeByUserID", mock.Anything, userId).
					Return(nil)
			},
			setupCtx: func(ctx context.Context) context.Context {
				return identity.WithPrincipal(ctx, &identity.Principal{
					UserID: userId,
					Role:   identity.RoleAdmin,
					Source: identity.AuthExternal,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshTokenRepo := &mocks.MockRefreshTokenRepository{}

			logoutMocks := &logoutMocks{
				refreshTokenRepo: refreshTokenRepo,
			}

			if tt.setupMocks != nil {
				tt.setupMocks(logoutMocks)
			}

			svcDeps := &auth.ServiceDeps{
				Logger:           slog.Default(),
				RefreshTokenRepo: refreshTokenRepo,
			}
			service := auth.NewService(svcDeps)

			ctx := t.Context()
			if tt.setupCtx != nil {
				ctx = tt.setupCtx(ctx)
			}

			err := service.LogoutAll(ctx)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
			} else {
				assert.NoError(t, err)
			}

			refreshTokenRepo.AssertExpectations(t)
		})
	}
}

func TestService_Refresh(t *testing.T) {
	t.Parallel()

	fakeHash := testdata.TestPasswordHash

	fakeRefreshToken := "fakeRefresh"
	defaultInput := auth.RefreshInput{
		RefreshToken: fakeRefreshToken,
	}

	defaultJti := uuid.New()
	defaultPrincipal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleUser,
		Source: identity.AuthExternal,
		JTI:    &defaultJti,
	}

	defaultUserReturn := &user.User{
		ID:       1,
		Name:     "test",
		UserName: "test",
		Email:    "test@test.com",
		Password: fakeHash,
		Role:     identity.RoleUser,
	}

	fakeAccessToken := "fakeAccess"
	fakeAccessTokenResult := &jwt.AccessTokenResult{
		Token: fakeAccessToken,
		Meta: jwt.TokenMetadata{
			JTI: defaultJti,
		},
	}

	tests := []struct {
		name       string
		input      auth.RefreshInput
		setupMocks func(m *refreshMocks)
		expected   *auth.RefreshOutput
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "bad refresh token",
			input: auth.RefreshInput{
				RefreshToken: "badFakeRefresh",
			},
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", "badFakeRefresh").
					Return(nil, errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "missing jti",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(&identity.Principal{UserID: 1}))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "db error",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Err[*user.User](gorm.ErrInvalidDB))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "non existing user",
			input: auth.RefreshInput{
				RefreshToken: "fakeRefreshTokenNonExistentUser",
			},
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", "fakeRefreshTokenNonExistentUser").
					Return(testutil.Ok(&identity.Principal{
						UserID: 999,
						Role:   identity.RoleUser,
						Source: identity.AuthExternal,
						JTI:    &defaultJti,
					}))

				m.userRepo.
					On("GetByID", mock.Anything, uint(999)).
					Return(testutil.Err[*user.User](gorm.ErrRecordNotFound))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "refresh token db error",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Err[*auth.RefreshToken](gorm.ErrInvalidDB))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "non existing refresh token",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Err[*auth.RefreshToken](gorm.ErrRecordNotFound))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "revoked refresh token",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Ok(&auth.RefreshToken{RevokedAt: testutil.Ptr(time.Now())}))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "expired refresh token",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Ok(&auth.RefreshToken{ExpiresAt: time.Now().Add(-1 * time.Hour)}))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "access token error",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Ok(&auth.RefreshToken{ExpiresAt: time.Now().Add(time.Hour)}))

				m.jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role, mock.Anything).
					Return(nil, errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "success",
			input: defaultInput,
			setupMocks: func(m *refreshMocks) {
				m.jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				m.userRepo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				m.refreshTokenRepo.
					On("GetByJTI", mock.Anything, mock.AnythingOfType("uuid.UUID")).
					Return(testutil.Ok(&auth.RefreshToken{ExpiresAt: time.Now().Add(time.Hour)}))

				m.jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role, mock.Anything).
					Return(fakeAccessTokenResult, nil)
			},
			expected: &auth.RefreshOutput{
				AccessToken: fakeAccessToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userRepo := &mocks.MockUserRepository{}
			jwtManager := &mocks.MockJwtManager{}
			refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
			refreshMocks := &refreshMocks{
				userRepo:         userRepo,
				jwtManager:       jwtManager,
				refreshTokenRepo: refreshTokenRepo,
			}
			tt.setupMocks(refreshMocks)

			svcDeps := &auth.ServiceDeps{
				AuthConfig: &config.AuthConfig{
					FakeHash: fakeHash,
				},
				UserRepo:         userRepo,
				RefreshTokenRepo: refreshTokenRepo,
				Logger:           slog.Default(),
				TokenManager:     jwtManager,
			}
			service := auth.NewService(svcDeps)

			result, err := service.Refresh(t.Context(), tt.input)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			userRepo.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
			jwtManager.AssertExpectations(t)
		})
	}
}
