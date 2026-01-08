package jwt

import (
	"errors"
	"gomonitor/internal/config"
	"gomonitor/internal/pkg/identity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidSignMethod = errors.New("invalid signing method")
	ErrInvalidTokenType  = errors.New("invalid token type")
)

type TokenManager struct {
	cfg *config.AuthConfig
}

func NewTokenManager(cfg *config.AuthConfig) *TokenManager {
	return &TokenManager{
		cfg: cfg,
	}
}

func (t *TokenManager) GenerateRefreshToken(userID uint, role identity.UserRole) (string, error) {
	return t.generateToken(userID, role, TokenTypeRefresh, t.cfg.RefreshTokenTTL, t.cfg.RefreshTokenSecret)
}

func (t *TokenManager) GenerateAccessToken(userID uint, role identity.UserRole) (string, error) {
	return t.generateToken(userID, role, TokenTypeAccess, t.cfg.AccessTokenTTL, t.cfg.AccessTokenSecret)
}

func (t *TokenManager) generateToken(userID uint, role identity.UserRole, tokenType TokenType, ttl time.Duration, secret string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		Type:   tokenType,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func (t *TokenManager) ValidateRefreshToken(tokenString string) (*identity.Principal, error) {
	return t.validateToken(tokenString, TokenTypeRefresh, t.cfg.RefreshTokenSecret)
}

func (t *TokenManager) ValidateAccessToken(tokenString string) (*identity.Principal, error) {
	return t.validateToken(tokenString, TokenTypeAccess, t.cfg.AccessTokenSecret)
}

func (t *TokenManager) validateToken(tokenString string, tokenType TokenType, secret string) (*identity.Principal, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.Type != tokenType {
		return nil, ErrInvalidTokenType
	}

	return &identity.Principal{
		UserID: claims.UserID,
		Role:   claims.Role,
		Source: identity.AuthExternal,
	}, nil
}
