package user

import (
	"gomonitor/internal/pkg/identity"
	"time"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	UserName  string
	Email     string            `gorm:"type:varchar(254);not null;uniqueIndex"`
	Password  string            `gorm:"type:char(60);not null"`
	Role      identity.UserRole `gorm:"type:user_role;not null;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
