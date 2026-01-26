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
