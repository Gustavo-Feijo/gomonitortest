package user_test

import (
	"context"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/domain/user/testdata"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Count(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		expectedResult testutil.Result[int64]
		expectError    bool
		setupFunc      func(db *gorm.DB)
		contextSetup   func(context.Context) context.Context
	}{
		{
			name: "number of seeded users matches count",
			expectedResult: testutil.Result[int64]{
				Value: 2,
				Err:   nil,
			},
			setupFunc: func(db *gorm.DB) {
				testdata.SeedUsers(t, db, 2)
			},
		},
		{
			name:           "fails if context cancelled",
			expectedResult: testutil.Result[int64]{Value: 0},
			expectError:    true,
			contextSetup:   testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t)

			if tt.setupFunc != nil {
				tt.setupFunc(tx)
			}

			repository := user.NewRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}

			result, err := repository.Count(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Zero(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult.Value, result)
		})
	}
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB)
		userToCreate *user.User
		expectError  bool
		contextSetup func(context.Context) context.Context
	}{
		{
			name: "successfully creates a user",
			userToCreate: &user.User{
				Name:     "John",
				UserName: "john",
				Email:    "john@test.com",
				Password: testdata.TestPasswordHash,
				Role:     identity.RoleUser,
			},
		},
		{
			name: "fails if context cancelled",
			userToCreate: &user.User{
				Name:     "doe",
				UserName: "doe",
				Email:    "doe@test.com",
				Password: testdata.TestPasswordHash,
				Role:     identity.RoleUser,
			},
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t)

			if tt.setupFunc != nil {
				tt.setupFunc(tx)
			}

			repo := user.NewRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}
			err := repo.Create(ctx, tt.userToCreate)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			var got user.User
			err = tx.WithContext(t.Context()).First(&got, "email = ?", tt.userToCreate.Email).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.userToCreate.UserName, got.UserName)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB) *user.User
		idToFind     uint
		expectError  bool
		contextSetup func(ctx context.Context) context.Context
	}{
		{
			name: "finds existing user",
			setupFunc: func(db *gorm.DB) *user.User {
				u := testdata.SeedUser(t, db, 0)
				return u
			},
			idToFind: 1,
		},
		{
			name:        "returns error if user does not exist",
			idToFind:    999,
			expectError: true,
		},
		{
			name:         "fails if context cancelled",
			idToFind:     1,
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t)

			var seededUser *user.User
			if tt.setupFunc != nil {
				seededUser = tt.setupFunc(tx)
				if seededUser != nil {
					tt.idToFind = seededUser.ID
				}
			}

			repo := user.NewRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}
			got, err := repo.GetByID(ctx, tt.idToFind)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, seededUser.ID, got.ID)
			assert.Equal(t, seededUser.Email, got.Email)
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB) *user.User
		emailToFind  string
		expectError  bool
		contextSetup func(ctx context.Context) context.Context
	}{
		{
			name: "finds existing user",
			setupFunc: func(db *gorm.DB) *user.User {
				u := testdata.SeedUser(t, db, 0)
				return u
			},
			emailToFind: "willbereplacedbyseed@test.com",
		},
		{
			name:        "returns error if user does not exist",
			emailToFind: "nonexistent@test.com",
			expectError: true,
		},
		{
			name:         "fails if context cancelled",
			emailToFind:  "willfail@test.com",
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t)

			var seededUser *user.User
			if tt.setupFunc != nil {
				seededUser = tt.setupFunc(tx)
				if seededUser != nil {
					tt.emailToFind = seededUser.Email
				}
			}

			repo := user.NewRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}
			got, err := repo.GetByEmail(ctx, tt.emailToFind)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, seededUser.ID, got.ID)
			assert.Equal(t, seededUser.Email, got.Email)
		})
	}
}
