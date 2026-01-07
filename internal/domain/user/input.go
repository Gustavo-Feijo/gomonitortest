package user

import "gomonitor/internal/pkg/identity"

type CreateUserInput struct {
	Name     string
	Email    string
	UserName string
	Password string
	Role     *identity.UserRole
}

type GetUserInput struct {
	ID uint
}
