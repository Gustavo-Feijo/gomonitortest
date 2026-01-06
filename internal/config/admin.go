package config

import "errors"

// App admin configuration.
type AdminConfig struct {
	Email    string
	Password string
}

func getAdminConfig() (*AdminConfig, error) {
	email := getEnv("ADMIN_EMAIL", "")
	password := getEnv("ADMIN_PASSWORD", "")

	if email == "" || password == "" {
		return nil, errors.New("missing admin configuration")
	}

	return &AdminConfig{
		Email:    email,
		Password: password,
	}, nil
}
