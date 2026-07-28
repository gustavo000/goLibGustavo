package telemetry

import (
	"os"
	"strconv"
	"strings"
)

// EnvKeys maps configuration fields to environment variable names.
// Override any field to use project-specific env vars in consuming services.
type EnvKeys struct {
	Enabled        string
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	OTLPProtocol   string
	Insecure       string
	SampleRatio    string
}

// DefaultEnvKeys follows OpenTelemetry environment variable conventions.
var DefaultEnvKeys = EnvKeys{
	Enabled:        "OTEL_ENABLED",
	ServiceName:    "OTEL_SERVICE_NAME",
	ServiceVersion: "OTEL_SERVICE_VERSION",
	Environment:    "OTEL_DEPLOYMENT_ENVIRONMENT",
	OTLPEndpoint:   "OTEL_EXPORTER_OTLP_ENDPOINT",
	OTLPProtocol:   "OTEL_EXPORTER_OTLP_PROTOCOL",
	Insecure:       "OTEL_EXPORTER_OTLP_INSECURE",
	SampleRatio:    "OTEL_TRACES_SAMPLER_ARG",
}

// Config holds telemetry provider settings.
type Config struct {
	Enabled        bool
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	OTLPProtocol   string
	Insecure       bool
	SampleRatio    float64
}

// LoadConfigFromEnv builds a Config using the provided env var mapping.
// When keys is nil, DefaultEnvKeys is used.
func LoadConfigFromEnv(keys *EnvKeys) Config {
	if keys == nil {
		keys = &DefaultEnvKeys
	}

	cfg := Config{
		Enabled:        envBool(keys.Enabled, true),
		ServiceName:    envString(keys.ServiceName, "unknown-service"),
		ServiceVersion: envString(keys.ServiceVersion, ""),
		Environment:    envString(keys.Environment, ""),
		OTLPEndpoint:   envString(keys.OTLPEndpoint, ""),
		OTLPProtocol:   normalizeProtocol(envString(keys.OTLPProtocol, "grpc")),
		Insecure:       envBool(keys.Insecure, false),
		SampleRatio:    envFloat(keys.SampleRatio, 1.0),
	}

	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"); endpoint != "" && cfg.OTLPEndpoint == "" {
		cfg.OTLPEndpoint = endpoint
	}
	if os.Getenv("OTEL_SDK_DISABLED") == "true" {
		cfg.Enabled = false
	}

	return cfg
}

func envString(key, fallback string) string {
	if key == "" {
		return fallback
	}
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if key == "" {
		return fallback
	}
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envFloat(key string, fallback float64) float64 {
	if key == "" {
		return fallback
	}
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func normalizeProtocol(protocol string) string {
	switch strings.ToLower(strings.TrimSpace(protocol)) {
	case "http", "http/protobuf", "http/protobufs":
		return "http"
	default:
		return "grpc"
	}
}
