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

type service struct {
	logger *slog.Logger
	repo   *repository
}

func NewService(deps *ServiceDeps) *service {
	repo := NewRepository(deps.DB)

	return &service{
		logger: deps.Logger,
		repo:   repo,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	principal, ok := identity.PrincipalFromContext(ctx)
	if !ok {
		return nil, errors.ErrUnauthenticated
	}

	if principal.Role != identity.RoleAdmin {
		return nil, errors.ErrForbidden
	}

	var role identity.UserRole
	if req.Role != nil {
		role = *req.Role
	} else {
		role = identity.RoleUser
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     req.Name,
		UserName: req.UserName,
		Email:    req.Email,
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
		"target_email", req.Email,
		"source", principal.Source,
	)

	resp := &CreateUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,

		CreatedAt: user.CreatedAt,
		UserName:  user.UserName,
	}

	return resp, nil
}

func (s *service) GetUser(ctx context.Context, id uint) (*GetUserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &GetUserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		UserName: user.UserName,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return resp, nil
}
