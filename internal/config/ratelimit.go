package config

import (
	"fmt"
	"time"
)

// RateLimit configuration.
type RateLimitConfig struct {
	IPLimit    int
	IPWindow   time.Duration
	UserLimit  int
	UserWindow time.Duration
}

func getRateLimitConfig() (*RateLimitConfig, error) {
	ipWindow := getEnv("RATE_LIMIT_IP_WINDOW", "1m")
	userWindow := getEnv("RATE_LIMIT_USER_WINDOW", "1m")

	ipWindowDuration, err := time.ParseDuration(ipWindow)
	if err != nil {
		return nil, fmt.Errorf("error parsing IpWindow: %v", err)
	}

	userWindowDuration, err := time.ParseDuration(userWindow)
	if err != nil {
		return nil, fmt.Errorf("error parsing UserWindow: %v", err)
	}

	ipLimit := getIntEnv("RATE_LIMIT_IP_LIMIT", 20)
	userLimit := getIntEnv("RATE_LIMIT_USER_LIMIT", 10)

	return &RateLimitConfig{
		IPLimit:    ipLimit,
		IPWindow:   ipWindowDuration,
		UserLimit:  userLimit,
		UserWindow: userWindowDuration,
	}, nil
}
