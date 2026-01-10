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
	Hasher   password.PasswordHasher
	Logger   *slog.Logger
	UserRepo Repository
}

type Service struct {
	hasher   password.PasswordHasher
	logger   *slog.Logger
	userRepo Repository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		logger:   deps.Logger,
		hasher:   deps.Hasher,
		userRepo: deps.UserRepo,
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

	hashedPassword, err := s.hasher.HashPassword(input.Password)
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

	err = s.userRepo.Create(ctx, user)
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
	user, err := s.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.NewNotFoundError("User not found", err)
		}
		return nil, err
	}

	return user, nil
}
