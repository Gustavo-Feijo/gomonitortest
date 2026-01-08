package auth

import (
	"context"
	"gomonitor/internal/config"
	"gomonitor/internal/observability/logging"
	pkgerrors "gomonitor/internal/pkg/errors"
	"gomonitor/internal/pkg/jwt"
	"gomonitor/internal/pkg/password"
	"log/slog"

	"gorm.io/gorm"
)

type ServiceDeps struct {
	AuthConfig   *config.AuthConfig
	DB           *gorm.DB
	Logger       *slog.Logger
	TokenManager *jwt.TokenManager
}

type Service struct {
	authCfg      *config.AuthConfig
	logger       *slog.Logger
	repo         Repository
	tokenManager *jwt.TokenManager
}

func NewService(deps *ServiceDeps) *Service {
	repo := NewRepository(deps.DB)

	return &Service{
		authCfg:      deps.AuthConfig,
		logger:       deps.Logger,
		repo:         repo,
		tokenManager: deps.TokenManager,
	}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	hash := s.authCfg.FakeHash
	if err == nil && user != nil {
		hash = user.Password
	}

	verifyErr := password.VerifyPassword(hash, input.Password)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	if user == nil || verifyErr != nil {
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

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	token, err := s.tokenManager.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		logging.FromContext(ctx).Warn(
			"unauthorized refresh request",
		)

		return nil, pkgerrors.NewUnauthorizedError(MsgInvalidToken)
	}

	user, err := s.repo.GetUserByID(ctx, token.UserID)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	if user == nil {
		logging.FromContext(ctx).Warn(
			"non existent user refresh request",
		)

		return nil, pkgerrors.NewUnauthorizedError(MsgInvalidToken)
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, pkgerrors.NewInternalError(err)
	}

	return &RefreshOutput{
		AccessToken: accessToken,
	}, nil
}
