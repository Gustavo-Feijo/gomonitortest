package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"gomonitor/internal/config"
	"gomonitor/internal/pkg/identity"
	pkgjwt "gomonitor/internal/pkg/jwt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testConfig = &config.AuthConfig{
	AccessTokenSecret:  "access",
	AccessTokenTTL:     time.Hour,
	RefreshTokenSecret: "refresh",
	RefreshTokenTTL:    time.Hour * 24,
}

func TestNewTokenManager(t *testing.T) {
	tokenManager := pkgjwt.NewTokenManager(&config.AuthConfig{})
	assert.NotNil(t, tokenManager)
}

func TestGenerateTokens(t *testing.T) {
	type tokenTestResult struct {
		Token     string
		JTI       uuid.UUID
		IssuedAt  time.Time
		ExpiresAt time.Time
	}

	tests := []struct {
		name      string
		generate  func(pkgjwt.TokenManager) (tokenTestResult, error)
		secret    string
		tokenType pkgjwt.TokenType
	}{
		{
			name: "access token",
			generate: func(tm pkgjwt.TokenManager) (tokenTestResult, error) {
				res, err := tm.GenerateAccessToken(1, identity.RoleAdmin)
				if err != nil {
					return tokenTestResult{}, err
				}

				return tokenTestResult{
					Token:     res.Token,
					IssuedAt:  res.Meta.IssuedAt,
					ExpiresAt: res.Meta.ExpiresAt,
				}, nil
			},
			secret:    testConfig.AccessTokenSecret,
			tokenType: pkgjwt.TokenTypeAccess,
		},
		{
			name: "refresh token",
			generate: func(tm pkgjwt.TokenManager) (tokenTestResult, error) {
				res, err := tm.GenerateRefreshToken(1, identity.RoleAdmin)
				if err != nil {
					return tokenTestResult{}, err
				}

				return tokenTestResult{
					Token:     res.Token,
					JTI:       res.Meta.JTI,
					IssuedAt:  res.Meta.IssuedAt,
					ExpiresAt: res.Meta.ExpiresAt,
				}, nil
			},
			secret:    testConfig.RefreshTokenSecret,
			tokenType: pkgjwt.TokenTypeRefresh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := pkgjwt.NewTokenManager(testConfig)

			tokenTestResult, err := tt.generate(tm)
			require.NoError(t, err)

			if tt.tokenType == pkgjwt.TokenTypeRefresh {
				assert.NotNil(t, tokenTestResult.JTI)
			}

			token, err := jwt.ParseWithClaims(
				tokenTestResult.Token,
				&pkgjwt.CustomClaims{},
				func(token *jwt.Token) (any, error) {
					return []byte(tt.secret), nil
				},
			)
			require.NoError(t, err)

			claims := token.Claims.(*pkgjwt.CustomClaims)
			assert.Equal(t, tt.tokenType, claims.Type)
		})
	}
}

func TestValidateRefreshToken(t *testing.T) {
	expectedPrincipal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthExternal,
	}

	tests := []struct {
		name        string
		tokenGen    func() string
		expectedErr error
	}{
		{
			name: "invalid token",
			tokenGen: func() string {
				return "invalid-jwt"
			},
			expectedErr: pkgjwt.ErrInvalidToken,
		},
		{
			name: "wrong secret",
			tokenGen: func() string {
				tm := pkgjwt.NewTokenManager(testConfig)
				token, _ := tm.GenerateAccessToken(1, identity.RoleUser)
				return token.Token
			},
			expectedErr: pkgjwt.ErrInvalidToken,
		},
		{
			name: "non-hmac algorithm",
			tokenGen: func() string {
				now := time.Now()
				claims := pkgjwt.CustomClaims{
					Type:   pkgjwt.TokenTypeRefresh,
					UserID: 1,
					Role:   identity.RoleAdmin,
					JTI:    uuid.New().String(),
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(testConfig.RefreshTokenTTL)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

				privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
				require.NoError(t, err)

				tokenStr, _ := token.SignedString(privateKey)
				return tokenStr

			},
			expectedErr: pkgjwt.ErrInvalidToken,
		},
		{
			name: "wrong access type",
			tokenGen: func() string {
				now := time.Now()
				claims := pkgjwt.CustomClaims{
					Type:   pkgjwt.TokenTypeAccess,
					UserID: 1,
					Role:   identity.RoleAdmin,
					JTI:    uuid.New().String(),
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(testConfig.RefreshTokenTTL)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, _ := token.SignedString([]byte(testConfig.RefreshTokenSecret))
				return tokenStr
			},
			expectedErr: pkgjwt.ErrInvalidTokenType,
		},
		{
			name: "success",
			tokenGen: func() string {
				now := time.Now()
				claims := pkgjwt.CustomClaims{
					Type:   pkgjwt.TokenTypeRefresh,
					UserID: 1,
					Role:   identity.RoleAdmin,
					JTI:    uuid.New().String(),
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(testConfig.RefreshTokenTTL)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, _ := token.SignedString([]byte(testConfig.RefreshTokenSecret))
				return tokenStr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := pkgjwt.NewTokenManager(testConfig)

			token := tt.tokenGen()

			principal, err := tm.ValidateRefreshToken(token)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, expectedPrincipal, principal)
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	expectedPrincipal := &identity.Principal{
		UserID: 1,
		Role:   identity.RoleAdmin,
		Source: identity.AuthExternal,
	}

	// Just success, logic always should be same as the refresh token generation, just changing type.
	tests := []struct {
		name        string
		tokenGen    func() string
		expectedErr error
	}{
		{
			name: "success",
			tokenGen: func() string {
				now := time.Now()
				claims := pkgjwt.CustomClaims{
					Type:   pkgjwt.TokenTypeAccess,
					UserID: 1,
					Role:   identity.RoleAdmin,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(now.Add(testConfig.AccessTokenTTL)),
						IssuedAt:  jwt.NewNumericDate(now),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				tokenStr, _ := token.SignedString([]byte(testConfig.AccessTokenSecret))
				return tokenStr
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := pkgjwt.NewTokenManager(testConfig)

			token := tt.tokenGen()

			principal, err := tm.ValidateAccessToken(token)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, tt.expectedErr, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, expectedPrincipal, principal)
			}
		})
	}
}
