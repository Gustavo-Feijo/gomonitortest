package config

// HTTP server configuration.
type HTTPConfig struct {
	Address string
}

// getHTTPConfig gets all necessary environment variables.
func getHTTPConfig() *HTTPConfig {
	port := getEnv("APP_PORT", "8080")
	return &HTTPConfig{
		Address: ":" + port,
	}
}
