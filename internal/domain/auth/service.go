package auth

import (
	"context"
	"errors"
	"gomonitor/internal/config"
	"gomonitor/internal/domain/user"
	"gomonitor/internal/observability/logging"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/jwt"
	"gomonitor/internal/pkg/password"
	"log/slog"

	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error)
}

type ServiceDeps struct {
	AuthConfig   *config.AuthConfig
	UserRepo     user.Repository
	Logger       *slog.Logger
	Hasher       password.PasswordHasher
	TokenManager jwt.TokenManager
}

type service struct {
	authCfg      *config.AuthConfig
	logger       *slog.Logger
	hasher       password.PasswordHasher
	userRepo     user.Repository
	tokenManager jwt.TokenManager
}

func NewService(deps *ServiceDeps) Service {
	return &service{
		authCfg:      deps.AuthConfig,
		logger:       deps.Logger,
		hasher:       deps.Hasher,
		userRepo:     deps.UserRepo,
		tokenManager: deps.TokenManager,
	}
}

func (s *service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	hash := s.authCfg.FakeHash
	if err == nil && user != nil {
		hash = user.Password
	}

	// First apply the hash to avoid enumeration.
	verifyErr := s.hasher.VerifyPassword(hash, input.Password)

	// Treat DB error or non existent user.
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.NewUnauthorizedError(MsgInvalidCredentials)
		}
		return nil, pkgerrors.NewInternalError(err)
	}

	if verifyErr != nil {
		logging.FromContext(ctx).Warn(
			"unauthorized login request",
			slog.Any("email", input.Email),
		)

		return nil, pkgerrors.NewUnauthorizedError(MsgInvalidCredentials)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	return &LoginOutput{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *service) Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	token, err := s.tokenManager.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		logging.FromContext(ctx).Warn(
			"unauthorized refresh request",
		)

		return nil, pkgerrors.NewUnauthorizedError(MsgInvalidToken)
	}

	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgerrors.NewUnauthorizedError(MsgInvalidCredentials)
		}
		return nil, pkgerrors.NewInternalError(err)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	return &RefreshOutput{
		AccessToken: accessToken,
	}, nil
}
