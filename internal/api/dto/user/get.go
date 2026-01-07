package userdto

import (
	"gomonitor/internal/domain/user"
	"gomonitor/internal/pkg/identity"
	"time"
)

type GetUserRequest struct {
	ID uint `uri:"id" binding:"required"`
}

func (r *GetUserRequest) ToDomainInput() user.GetUserInput {
	return user.GetUserInput{
		ID: r.ID,
	}
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

func ToGetUserResponse(user *user.User) *GetUserResponse {
	return &GetUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		UserName: user.UserName,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
