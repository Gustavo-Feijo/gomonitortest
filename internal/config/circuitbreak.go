package config

// Circuit Breaker configuration.
type CircuitBreakerConfig struct {
	MaxRequests uint32
	MaxFailures uint32
}

func getCircuitBreakerConfig() *CircuitBreakerConfig {
	maxRequests := getIntEnv("CIRCUIT_BREAKER_MAX_REQUEST", 5)
	maxFailures := getIntEnv("CIRCUIT_BREAKER_MAX_FAILURES", 5)

	return &CircuitBreakerConfig{
		MaxRequests: uint32(maxRequests),
		MaxFailures: uint32(maxFailures),
	}
}
