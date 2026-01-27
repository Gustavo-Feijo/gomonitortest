package jwt

import (
	"errors"
	"gomonitor/internal/config"
	"gomonitor/internal/pkg/identity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidSignMethod = errors.New("invalid signing method")
	ErrInvalidTokenType  = errors.New("invalid token type")
)

type TokenManager interface {
	GenerateRefreshToken(userID uint, role identity.UserRole) (*RefreshTokenResult, error)
	GenerateAccessToken(userID uint, role identity.UserRole, refreshTokenJTI uuid.UUID) (*AccessTokenResult, error)
	ValidateRefreshToken(tokenString string) (*identity.Principal, error)
	ValidateAccessToken(tokenString string) (*identity.Principal, error)
}

type tokenManager struct {
	cfg *config.AuthConfig
}

func NewTokenManager(cfg *config.AuthConfig) TokenManager {
	return &tokenManager{
		cfg: cfg,
	}
}

type RefreshTokenResult struct {
	Token string
	Meta  TokenMetadata
}

type AccessTokenResult struct {
	Token string
	Meta  TokenMetadata
}

type TokenMetadata struct {
	JTI       uuid.UUID
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func (t *tokenManager) GenerateRefreshToken(userID uint, role identity.UserRole) (*RefreshTokenResult, error) {
	now := time.Now()
	expiresAt := now.Add(t.cfg.RefreshTokenTTL)
	token, metadata, err := t.generateToken(userID, role, TokenTypeRefresh, uuid.New(), expiresAt, now, t.cfg.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		Token: token,
		Meta:  metadata,
	}, nil
}

func (t *tokenManager) GenerateAccessToken(userID uint, role identity.UserRole, refreshTokenJTI uuid.UUID) (*AccessTokenResult, error) {
	now := time.Now()
	expiresAt := now.Add(t.cfg.AccessTokenTTL)
	token, metadata, err := t.generateToken(userID, role, TokenTypeAccess, refreshTokenJTI, expiresAt, now, t.cfg.AccessTokenSecret)
	if err != nil {
		return nil, err
	}

	return &AccessTokenResult{
		Token: token,
		Meta:  metadata,
	}, nil
}

func (t *tokenManager) generateToken(
	userID uint,
	role identity.UserRole,
	tokenType TokenType,
	jtiUUID uuid.UUID,
	expiresAt time.Time,
	issuedAt time.Time,
	secret string,
) (string, TokenMetadata, error) {
	claims := CustomClaims{
		Type:   tokenType,
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	if tokenType == TokenTypeRefresh {
		claims.JTI = jtiUUID.String()
	}

	if tokenType == TokenTypeAccess {
		claims.RefreshJTI = jtiUUID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(secret))
	tokenMetadata := TokenMetadata{
		JTI:       jtiUUID,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}

	return tokenStr, tokenMetadata, err
}

func (t *tokenManager) ValidateRefreshToken(tokenString string) (*identity.Principal, error) {
	return t.validateToken(tokenString, TokenTypeRefresh, t.cfg.RefreshTokenSecret)
}

func (t *tokenManager) ValidateAccessToken(tokenString string) (*identity.Principal, error) {
	return t.validateToken(tokenString, TokenTypeAccess, t.cfg.AccessTokenSecret)
}

func (t *tokenManager) validateToken(tokenString string, tokenType TokenType, secret string) (*identity.Principal, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignMethod
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Guaranteed to be *CustomClains due to parse with Claims.
	claims, _ := token.Claims.(*CustomClaims)
	if claims.Type != tokenType {
		return nil, ErrInvalidTokenType
	}

	var jti *uuid.UUID
	if tokenType == TokenTypeRefresh {
		parsed, err := uuid.Parse(claims.JTI)
		if err != nil {
			return nil, ErrInvalidToken
		}
		jti = &parsed
	}

	var refreshjti *uuid.UUID
	if tokenType == TokenTypeAccess {
		parsed, err := uuid.Parse(claims.RefreshJTI)
		if err != nil {
			return nil, ErrInvalidToken
		}
		refreshjti = &parsed
	}

	return &identity.Principal{
		UserID:     claims.UserID,
		Role:       claims.Role,
		Source:     identity.AuthExternal,
		JTI:        jti,
		RefreshJTI: refreshjti,
	}, nil
}
