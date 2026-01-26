package user_test

import (
	"context"
	"errors"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/infra/database/postgres"
	"gomonitor/internal/mocks"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/password"
	"gomonitor/internal/testutil"
	"log/slog"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	deps := &user.ServiceDeps{
		Hasher:   password.NewPasswordHasher(bcrypt.DefaultCost),
		UserRepo: user.NewUserRepository(&gorm.DB{}),
		Logger:   &slog.Logger{},
	}

	service := user.NewService(deps)
	assert.NotNil(t, service)
}

// Match only the user email on the returned user, matching whole struct becomes cumbersome due to hashing and timestamps.
func matchUserEmail(email string) any {
	return mock.MatchedBy(func(u *user.User) bool {
		return u.Email == email
	})
}

func TestService_CreateUser(t *testing.T) {
	defaultInput := user.CreateUserInput{
		Name:     "test1",
		Email:    "test@test.com",
		UserName: "test1",
		Password: "password123",
	}

	newAdminInput := user.CreateUserInput{
		Name:     "TestAdmin",
		Email:    "admin@admin.com",
		UserName: "admin1",
		Password: "adminpswd",
		Role:     testutil.Ptr(identity.RoleAdmin),
	}

	adminCtx := func(ctx context.Context) context.Context {
		return identity.WithPrincipal(ctx, &identity.Principal{
			UserID: 1,
			Role:   identity.RoleAdmin,
			Source: identity.AuthInternal,
		})
	}

	userCtx := func(ctx context.Context) context.Context {
		return identity.WithPrincipal(ctx, &identity.Principal{
			UserID: 2,
			Role:   identity.RoleUser,
			Source: identity.AuthInternal,
		})
	}

	tests := []struct {
		name      string
		input     user.CreateUserInput
		setupMock func(repo *mocks.MockUserRepository)
		setupCtx  func(ctx context.Context) context.Context
		expected  *user.User
		assertErr func(t *testing.T, err error)
	}{
		{
			name:      "unauthorized due to missing principal",
			input:     defaultInput,
			setupMock: func(repo *mocks.MockUserRepository) {},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name:      "unauthorized due to user without admin role",
			input:     defaultInput,
			setupMock: func(repo *mocks.MockUserRepository) {},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
			setupCtx: userCtx,
		},
		{
			name: "password too long",
			input: user.CreateUserInput{
				Name:     "test1",
				Email:    "test@test.com",
				UserName: "test1",
				Password: strings.Repeat("a", 100),
				Role:     testutil.Ptr(identity.RoleUser),
			},
			setupMock: func(repo *mocks.MockUserRepository) {},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
			},
			setupCtx: adminCtx,
		},
		{
			name:  "database error duplicate constraint",
			input: defaultInput,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("Create", mock.Anything, matchUserEmail(defaultInput.Email)).
					Return(&pgconn.PgError{Code: postgres.UniqueViolation})
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
			setupCtx: adminCtx,
		},
		{
			name:  "random database error",
			input: defaultInput,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("Create", mock.Anything, matchUserEmail(defaultInput.Email)).
					Return(errors.New("DB error"))
			},
			assertErr: func(t *testing.T, err error) {
				assert.ErrorContains(t, err, "DB error")
			},
			setupCtx: adminCtx,
		},

		{
			name:  "success",
			input: newAdminInput,
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("Create", mock.Anything, matchUserEmail(newAdminInput.Email)).
					Return(nil)
			},
			setupCtx: adminCtx,
			expected: &user.User{Email: newAdminInput.Email},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.MockUserRepository{}
			tt.setupMock(repo)

			svcDeps := &user.ServiceDeps{
				Hasher:   password.NewPasswordHasher(bcrypt.DefaultCost),
				UserRepo: repo,
			}
			service := user.NewService(svcDeps)

			ctx := t.Context()
			if tt.setupCtx != nil {
				ctx = tt.setupCtx(ctx)
			}

			result, err := service.CreateUser(ctx, tt.input)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Email, result.Email)
			}
			repo.AssertExpectations(t)
		})
	}
}
func TestService_GetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     user.GetUserInput
		setupMock func(repo *mocks.MockUserRepository)
		expected  *user.User
		assertErr func(t *testing.T, err error)
	}{
		{
			name: "success",
			input: user.GetUserInput{
				ID: 1,
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("GetByID", mock.Anything, uint(1)).
					Return(testutil.Ok(&user.User{
						ID:    1,
						Email: "test@test.com"},
					))
			},
			expected: &user.User{ID: 1, Email: "test@test.com"},
		},
		{
			name: "user not found",
			input: user.GetUserInput{
				ID: 2,
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("GetByID", mock.Anything, uint(2)).
					Return(testutil.Err[*user.User](gorm.ErrRecordNotFound))
			},
			assertErr: func(t *testing.T, err error) {
				var nf *pkgerrors.AppError
				assert.ErrorAs(t, err, &nf)
			},
		},
		{
			name: "repository error",
			input: user.GetUserInput{
				ID: 3,
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.
					On("GetByID", mock.Anything, uint(3)).
					Return(nil, errors.New("db down"))
			},
			assertErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "db down")
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := &mocks.MockUserRepository{}
			tt.setupMock(repo)

			svcDeps := &user.ServiceDeps{
				Hasher:   password.NewPasswordHasher(bcrypt.DefaultCost),
				UserRepo: repo,
			}
			service := user.NewService(svcDeps)

			result, err := service.GetUser(t.Context(), tt.input)

			if tt.assertErr != nil {
				assert.Error(t, err)
				tt.assertErr(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			repo.AssertExpectations(t)
		})
	}
}
