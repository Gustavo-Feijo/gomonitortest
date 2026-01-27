package auth_test

import (
	"context"
	"gomonitor/internal/domain/auth"
	databaseinfra "gomonitor/internal/infra/database"
	"gomonitor/internal/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	t.Parallel()
	db, _ := databaseinfra.New(t.Context(), testDbCfg)

	tests := []struct {
		name                 string
		setupFunc            func(db *gorm.DB)
		refreshTokenToCreate *auth.RefreshToken
		expectError          bool
		contextSetup         func(context.Context) context.Context
	}{
		{
			name: "successfully creates a user",
			refreshTokenToCreate: &auth.RefreshToken{
				JTI:       uuid.New(),
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Hour * 24),
				CreatedAt: time.Now(),
			},
		},
		{
			name: "fails if context cancelled",
			refreshTokenToCreate: &auth.RefreshToken{
				JTI:       uuid.New(),
				UserID:    1,
				ExpiresAt: time.Now().Add(time.Hour * 24),
				CreatedAt: time.Now(),
			},
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t, db)

			if tt.setupFunc != nil {
				tt.setupFunc(tx)
			}

			repo := auth.NewRefreshTokenRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}
			err := repo.Create(ctx, tt.refreshTokenToCreate)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			var got auth.RefreshToken
			err = tx.WithContext(t.Context()).First(&got, "jti = ?", tt.refreshTokenToCreate.JTI).Error
			assert.NoError(t, err)
			assert.Equal(t, tt.refreshTokenToCreate.JTI, got.JTI)
		})
	}
}

func TestRepository_GetByJTI(t *testing.T) {
	t.Parallel()
	db, _ := databaseinfra.New(t.Context(), testDbCfg)

	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB, token *auth.RefreshToken)
		jti          uuid.UUID
		expectError  bool
		contextSetup func(context.Context) context.Context
	}{
		{
			name: "successfully gets refresh token by jti",
			setupFunc: func(db *gorm.DB, token *auth.RefreshToken) {
				err := db.Create(token).Error
				assert.NoError(t, err)
			},
		},
		{
			name:        "returns error when token not found",
			jti:         uuid.New(),
			expectError: true,
		},
		{
			name: "fails if context cancelled",
			setupFunc: func(db *gorm.DB, token *auth.RefreshToken) {
				err := db.Create(token).Error
				assert.NoError(t, err)
			},
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t, db)

			token := &auth.RefreshToken{
				JTI:       uuid.New(),
				UserID:    1,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
			}

			if tt.jti == uuid.Nil {
				tt.jti = token.JTI
			}

			if tt.setupFunc != nil {
				tt.setupFunc(tx, token)
			}

			repo := auth.NewRefreshTokenRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}

			got, err := repo.GetByJTI(ctx, tt.jti)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, token.JTI, got.JTI)
		})
	}
}

func TestRepository_RevokeByJTI(t *testing.T) {
	t.Parallel()
	db, _ := databaseinfra.New(t.Context(), testDbCfg)

	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB, token *auth.RefreshToken)
		expectError  bool
		contextSetup func(context.Context) context.Context
	}{
		{
			name: "successfully revokes refresh token by jti",
			setupFunc: func(db *gorm.DB, token *auth.RefreshToken) {
				err := db.Create(token).Error
				assert.NoError(t, err)
			},
		},
		{
			name:         "fails if context cancelled",
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
			setupFunc: func(db *gorm.DB, token *auth.RefreshToken) {
				err := db.Create(token).Error
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t, db)

			token := &auth.RefreshToken{
				JTI:       uuid.New(),
				UserID:    1,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
			}

			if tt.setupFunc != nil {
				tt.setupFunc(tx, token)
			}

			repo := auth.NewRefreshTokenRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}

			err := repo.RevokeByJTI(ctx, token.JTI)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			var got auth.RefreshToken
			err = tx.First(&got, "jti = ?", token.JTI).Error
			assert.NoError(t, err)
			assert.NotNil(t, got.RevokedAt)
		})
	}
}

func TestRepository_RevokeByUserID(t *testing.T) {
	t.Parallel()
	db, _ := databaseinfra.New(t.Context(), testDbCfg)

	tests := []struct {
		name         string
		setupFunc    func(db *gorm.DB, userID uint)
		userID       uint
		expectError  bool
		contextSetup func(context.Context) context.Context
	}{
		{
			name:   "successfully revokes all refresh tokens for user",
			userID: 1,
			setupFunc: func(db *gorm.DB, userID uint) {
				tokens := []auth.RefreshToken{
					{
						JTI:       uuid.New(),
						UserID:    userID,
						ExpiresAt: time.Now().Add(24 * time.Hour),
						CreatedAt: time.Now(),
					},
					{
						JTI:       uuid.New(),
						UserID:    userID,
						ExpiresAt: time.Now().Add(24 * time.Hour),
						CreatedAt: time.Now(),
					},
				}

				err := db.Create(&tokens).Error
				assert.NoError(t, err)
			},
		},
		{
			name:         "fails if context cancelled",
			userID:       1,
			expectError:  true,
			contextSetup: testutil.GetCancelledCtx,
			setupFunc: func(db *gorm.DB, userID uint) {
				err := db.Create(&auth.RefreshToken{
					JTI:       uuid.New(),
					UserID:    userID,
					ExpiresAt: time.Now().Add(24 * time.Hour),
					CreatedAt: time.Now(),
				}).Error
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := setupTx(t, db)

			if tt.setupFunc != nil {
				tt.setupFunc(tx, tt.userID)
			}

			repo := auth.NewRefreshTokenRepository(tx)

			ctx := t.Context()
			if tt.contextSetup != nil {
				ctx = tt.contextSetup(ctx)
			}

			err := repo.RevokeByUserID(ctx, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			var count int64
			err = tx.
				Model(&auth.RefreshToken{}).
				Where("user_id = ? AND revoked_at IS NULL", tt.userID).
				Count(&count).
				Error

			assert.NoError(t, err)
			assert.Equal(t, int64(0), count)
		})
	}
}
