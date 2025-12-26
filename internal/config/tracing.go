package config

// Tracing configuration.
type TracingConfig struct {
	ServiceName string
	Address     string
}

func getTracingConfig() *TracingConfig {
	return &TracingConfig{
		ServiceName: getEnv("TRACE_SERVICE_NAME", "monitor"),
		Address:     getEnv("TRACE_SERVICE_ADDRESS", "tempo:4317"),
	}
}
