package config

// Logging configuration.
type LoggingConfig struct {
	Level string
}

func getLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Level: getEnv("LOG_LEVEL", "debug"),
	}
}
