package config

import (
	"errors"
	"fmt"
	"time"
)

// App auth configuration.
type AuthConfig struct {
	AccessTokenSecret  string
	AccessTokenTTL     time.Duration
	FakeHash           string
	RefreshTokenSecret string
	RefreshTokenTTL    time.Duration
}

func getAuthConfig() (*AuthConfig, error) {
	AccessTokenSecret := getEnv("AUTH_ACCESS_TOKEN_SECRET", "")
	refreshTokenSecret := getEnv("AUTH_REFRESH_TOKEN_SECRET", "")

	// Fake hash to use when user doesn't exist.
	fakeHash := getEnv("AUTH_FAKE_HASH", "")

	if AccessTokenSecret == "" {
		return nil, errors.New("missing AUTH_ACCESS_TOKEN_SECRET")
	}

	if refreshTokenSecret == "" {
		return nil, errors.New("missing AUTH_REFRESH_TOKEN_SECRET")
	}

	if fakeHash == "" {
		return nil, errors.New("missing AUTH_FAKE_HASH")
	}

	AccessTokenTTL := getEnv("AUTH_ACCESS_TOKEN_TTL", "1h")
	refreshTokenTTL := getEnv("AUTH_REFRESH_TOKEN_TTL", "168h")

	accessTokenDuration, err := time.ParseDuration(AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("error parsing AccessTokenTTL: %v", err)
	}

	refreshTokenDuration, err := time.ParseDuration(refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("error parsing refreshTokenTTL: %v", err)
	}

	return &AuthConfig{
		AccessTokenSecret:  AccessTokenSecret,
		AccessTokenTTL:     accessTokenDuration,
		FakeHash:           fakeHash,
		RefreshTokenSecret: refreshTokenSecret,
		RefreshTokenTTL:    refreshTokenDuration,
	}, nil
}
