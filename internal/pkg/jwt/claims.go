package jwt

import (
	"gomonitor/internal/pkg/identity"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type CustomClaims struct {
	Type   TokenType `json:"typ"`
	UserID uint      `json:"sub"`
	Role   identity.UserRole
	jwt.RegisteredClaims
}
