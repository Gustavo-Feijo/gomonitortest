package userdto

import (
	"gomonitor/internal/domain/user"
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

func (r *CreateUserRequest) ToDomainInput() user.CreateUserInput {
	return user.CreateUserInput{
		Name:     r.Name,
		Email:    r.Email,
		UserName: r.UserName,
		Password: r.Password,
		Role:     r.Role,
	}
}

type CreateUserResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	UserName  string            `json:"username"`
	Role      identity.UserRole `json:"role,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

func ToCreateUserResponse(user *user.User) *CreateUserResponse {
	return &CreateUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,

		CreatedAt: user.CreatedAt,
		UserName:  user.UserName,
	}
}
