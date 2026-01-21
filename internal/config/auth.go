package config

import (
	"fmt"
	"strings"
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
	var missing []string
	accessToken := getEnv("AUTH_ACCESS_TOKEN_SECRET", "")
	refreshToken := getEnv("AUTH_REFRESH_TOKEN_SECRET", "")
	fakeHash := getEnv("AUTH_FAKE_HASH", "")

	if accessToken == "" {
		missing = append(missing, "AUTH_ACCESS_TOKEN_SECRET")
	}
	if refreshToken == "" {
		missing = append(missing, "AUTH_REFRESH_TOKEN_SECRET")
	}
	if fakeHash == "" {
		missing = append(missing, "AUTH_FAKE_HASH")
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing auth config: %s", strings.Join(missing, ", "))
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
		AccessTokenSecret:  accessToken,
		AccessTokenTTL:     accessTokenDuration,
		FakeHash:           fakeHash,
		RefreshTokenSecret: refreshToken,
		RefreshTokenTTL:    refreshTokenDuration,
	}, nil
}
