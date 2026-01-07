package user

import (
	"context"
	"gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/password"
	"log/slog"

	"gorm.io/gorm"
)

type ServiceDeps struct {
	Logger *slog.Logger
	DB     *gorm.DB
}

type Service struct {
	logger *slog.Logger
	repo   *repository
}

func NewService(deps *ServiceDeps) *Service {
	repo := NewRepository(deps.DB)

	return &Service{
		logger: deps.Logger,
		repo:   repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return nil, errors.ErrUnauthenticated
	}

	if principal.Role != identity.RoleAdmin {
		return nil, errors.ErrForbidden
	}

	var role identity.UserRole
	if input.Role != nil {
		role = *input.Role
	} else {
		role = identity.RoleUser
	}

	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     input.Name,
		UserName: input.UserName,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     role,
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user",
			"created_by", principal.UserID,
			"email", user.Email,
			"error", err,
		)
		return nil, err
	}

	s.logger.Info("user created",
		"created_by", principal.UserID,
		"target_email", input.Email,
		"source", principal.Source,
	)

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, input GetUserInput) (*User, error) {
	user, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
