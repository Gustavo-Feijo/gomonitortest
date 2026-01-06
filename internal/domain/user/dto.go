package user

import (
	"gomonitor/internal/pkg/identity"
	"time"
)

type CreateUserRequest struct {
	Name     string             `json:"name" binding:"required"`
	Email    string             `json:"email" binding:"required,email"`
	UserName string             `json:"username" binding:"required"`
	Password string             `json:"password" binding:"required,min=8,max=72"`
	Role     *identity.UserRole `json:"role" binding:"omitempty,oneof=admin user"`
}

type CreateUserResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	UserName  string            `json:"username"`
	Role      identity.UserRole `json:"role,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

type GetUserRequest struct {
}

type GetUserResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	UserName  string            `json:"username"`
	Role      identity.UserRole `json:"role,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time
}
