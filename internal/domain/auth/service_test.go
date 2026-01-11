package auth_test

import (
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	deps := &auth.ServiceDeps{
		UserRepo:     user.NewRepository(&gorm.DB{}),
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

	tests := []struct {
		name       string
		input      auth.LoginInput
		setupMocks func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager)
		expected   *auth.LoginOutput
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "non existing user",
			input: auth.LoginInput{
				Email:    "nonexistenttest@test.com",
				Password: "password123",
			},
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "nonexistenttest@test.com").
					Return(testutil.Err[*user.User](gorm.ErrRecordNotFound))

				hasher.
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
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Err[*user.User](gorm.ErrInvalidDB))

				hasher.
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
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				hasher.
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
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return("", errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "access token error",
			input: defaultInput,
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeRefreshToken, nil)

				jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return("", errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "success",
			input: defaultInput,
			setupMocks: func(repo *mocks.MockUserRepository, hasher *mocks.MockPasswordHasher, jwtManager *mocks.MockJwtManager) {
				repo.
					On("GetByEmail", mock.Anything, "test@test.com").
					Return(testutil.Ok(defaultUserReturn))

				hasher.
					On("VerifyPassword", defaultUserReturn.Password, "password123").
					Return(nil)

				jwtManager.
					On("GenerateRefreshToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeRefreshToken, nil)

				jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeAccessToken, nil)
			},
			expected: &auth.LoginOutput{
				RefreshToken: fakeRefreshToken,
				AccessToken:  fakeAccessToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.MockUserRepository{}
			hasher := &mocks.MockPasswordHasher{}
			jwtManager := &mocks.MockJwtManager{}
			tt.setupMocks(repo, hasher, jwtManager)

			svcDeps := &auth.ServiceDeps{
				AuthConfig: &config.AuthConfig{
					FakeHash: fakeHash,
				},
				Hasher:       hasher,
				UserRepo:     repo,
				Logger:       slog.Default(),
				TokenManager: jwtManager,
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

			repo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			jwtManager.AssertExpectations(t)
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

	defaultPrincipal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleUser,
		Source: identity.AuthExternal,
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

	tests := []struct {
		name       string
		input      auth.RefreshInput
		setupMocks func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager)
		expected   *auth.RefreshOutput
		assertErr  func(t *testing.T, err error)
	}{
		{
			name: "bad refresh token",
			input: auth.RefreshInput{
				RefreshToken: "badFakeRefresh",
			},
			setupMocks: func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager) {
				jwtManager.
					On("ValidateRefreshToken", "badFakeRefresh").
					Return(nil, errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "db error",
			input: defaultInput,
			setupMocks: func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager) {
				jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				repo.
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
			setupMocks: func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager) {
				jwtManager.
					On("ValidateRefreshToken", "fakeRefreshTokenNonExistentUser").
					Return(testutil.Ok(&identity.Principal{
						UserID: 999,
						Role:   identity.RoleUser,
						Source: identity.AuthExternal,
					}))

				repo.
					On("GetByID", mock.Anything, uint(999)).
					Return(testutil.Err[*user.User](gorm.ErrRecordNotFound))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "access token error",
			input: defaultInput,
			setupMocks: func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager) {
				jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				repo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return("", errors.New("signing error"))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:  "success",
			input: defaultInput,
			setupMocks: func(repo *mocks.MockUserRepository, jwtManager *mocks.MockJwtManager) {
				jwtManager.
					On("ValidateRefreshToken", fakeRefreshToken).
					Return(testutil.Ok(defaultPrincipal))

				repo.
					On("GetByID", mock.Anything, defaultPrincipal.UserID).
					Return(testutil.Ok(defaultUserReturn))

				jwtManager.
					On("GenerateAccessToken", defaultUserReturn.ID, defaultUserReturn.Role).
					Return(fakeAccessToken, nil)
			},
			expected: &auth.RefreshOutput{
				AccessToken: fakeAccessToken,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.MockUserRepository{}
			jwtManager := &mocks.MockJwtManager{}
			tt.setupMocks(repo, jwtManager)

			svcDeps := &auth.ServiceDeps{
				AuthConfig: &config.AuthConfig{
					FakeHash: fakeHash,
				},
				UserRepo:     repo,
				Logger:       slog.Default(),
				TokenManager: jwtManager,
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

			repo.AssertExpectations(t)
			jwtManager.AssertExpectations(t)
		})
	}
}
