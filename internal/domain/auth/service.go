package auth

import (
	"context"
	"gomonitor/internal/config"
	"gomonitor/internal/pkg/password"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type ServiceDeps struct {
	AuthConfig *config.AuthConfig
	DB         *gorm.DB
	Logger     *slog.Logger
}

type Service struct {
	authCfg *config.AuthConfig
	logger  *slog.Logger
	repo    Repository
}

func NewService(deps *ServiceDeps) *Service {
	repo := NewRepository(deps.DB)

	return &Service{
		authCfg: deps.AuthConfig,
		logger:  deps.Logger,
		repo:    repo,
	}
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := s.repo.GetUserByEmail(ctx, input.Email)

	hash := s.authCfg.FakeHash
	if err == nil && user != nil {
		hash = user.Password
	}

	verifyErr := password.VerifyPassword(hash, input.Password)

	if err != nil || user == nil || verifyErr != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()

	refreshTokenClaim := jwt.MapClaims{
		"typ": "refresh",
		"sub": user.ID,
		// 7 days of expiration to the long lived refresh token.
		"exp": now.Add(time.Hour * 24 * 7).Unix(),
		"iat": now.Unix(),
	}

	apiTokenClaim := jwt.MapClaims{
		"typ": "access",
		"sub": user.ID,
		// 1 hour of expiration to the short lived API token.
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaim)
	apiToken := jwt.NewWithClaims(jwt.SigningMethodHS256, apiTokenClaim)

	refreshTokenStr, err := refreshToken.SignedString([]byte(s.authCfg.RefreshTokenSecret))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	apiTokenStr, err := apiToken.SignedString([]byte(s.authCfg.ApiTokenSecret))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LoginOutput{
		RefreshToken: refreshTokenStr,
		ApiToken:     apiTokenStr,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	token, err := jwt.Parse(input.RefreshToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.authCfg.RefreshTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
		return nil, ErrInvalidToken
	}

	sub, ok := claims["sub"]
	if !ok {
		return nil, ErrInvalidToken
	}

	subFloat, ok := sub.(float64)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID := uint(subFloat)

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrInvalidToken
	}

	now := time.Now()

	accessClaims := jwt.MapClaims{
		"typ": "access",
		"sub": user.ID,
		"exp": now.Add(time.Hour).Unix(),
		"iat": now.Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(s.authCfg.ApiTokenSecret))
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &RefreshOutput{
		ApiToken: accessToken,
	}, nil
}
