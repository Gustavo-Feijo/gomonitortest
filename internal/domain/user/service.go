package user

import (
	"context"
	"errors"
	"gomonitor/internal/infra/database/postgres"
	"gomonitor/internal/observability/logging"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/identity"
	"gomonitor/internal/pkg/password"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
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
		logging.FromContext(ctx).Warn("unauthenticated user creation attempt")
		return nil, pkgerrors.NewUnauthorizedError("unauthenticated")
	}

	if principal.Role != identity.RoleAdmin {
		logging.FromContext(ctx).Warn("unauthorized user creation attempt",
			"user_id", principal.UserID,
			"user_role", principal.Role,
			"source", principal.Source,
		)
		return nil, pkgerrors.NewForbiddenError()
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
		logging.FromContext(ctx).Error("failed to create user",
			"created_by", principal.UserID,
			"email", user.Email,
			"error", err,
		)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == postgres.UniqueViolation {
			return nil, pkgerrors.NewConflictError("Duplicate entry", err)
		}

		return nil, err
	}

	logging.FromContext(ctx).Info("user created",
		"created_by", principal.UserID,
		"target_email", input.Email,
		"source", principal.Source,
	)

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, input GetUserInput) (*User, error) {
	user, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pkgerrors.NewNotFoundError("User not found", err)
		}
		return nil, err
	}

	return user, nil
}
