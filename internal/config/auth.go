package config

import "errors"

// App auth configuration.
type AuthConfig struct {
	ApiTokenSecret     string
	FakeHash           string
	RefreshTokenSecret string
}

func getAuthConfig() (*AuthConfig, error) {
	apiTokenSecret := getEnv("AUTH_API_TOKEN_SECRET", "")
	refreshTokenSecret := getEnv("AUTH_REFRESH_TOKEN_SECRET", "")

	// Fake hash to use when user doesn't exist.
	fakeHash := getEnv("AUTH_FAKE_HASH", "")

	if apiTokenSecret == "" {
		return nil, errors.New("missing AUTH_API_TOKEN_SECRET")
	}

	if refreshTokenSecret == "" {
		return nil, errors.New("missing AUTH_REFRESH_TOKEN_SECRET")
	}

	if fakeHash == "" {
		return nil, errors.New("missing AUTH_FAKE_HASH")
	}

	return &AuthConfig{
		ApiTokenSecret:     apiTokenSecret,
		FakeHash:           fakeHash,
		RefreshTokenSecret: refreshTokenSecret,
	}, nil
}
