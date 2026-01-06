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

type service struct {
	authCfg *config.AuthConfig
	logger  *slog.Logger
	repo    Repository
}

func NewService(deps *ServiceDeps) *service {
	repo := NewRepository(deps.DB)

	return &service{
		authCfg: deps.AuthConfig,
		logger:  deps.Logger,
		repo:    repo,
	}
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)

	hash := s.authCfg.FakeHash
	if err == nil && user != nil {
		hash = user.Password
	}

	verifyErr := password.VerifyPassword(hash, req.Password)

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

	return &LoginResponse{
		RefreshToken: refreshTokenStr,
		ApiToken:     apiTokenStr,
	}, nil
}

func (s *service) Refresh(ctx context.Context, req RefreshRequest) (*RefreshResponse, error) {
	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (any, error) {
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

	return &RefreshResponse{
		ApiToken: accessToken,
	}, nil
}
