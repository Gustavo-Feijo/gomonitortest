package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	JTI       uuid.UUID `gorm:"type:uuid;primaryKey;column:jti"`
	UserID    uint      `gorm:"index;column:user_id"`
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}
