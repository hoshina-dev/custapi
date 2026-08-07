package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hoshina-dev/custapi/internal/telemetry"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port           int
	DataSourceName string
	CorsOrigins    string
	Telemetry      telemetry.Config
}

// Load loads configuration from environment variables
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}
	port := getEnvInt("PORT", 8080)

	// Build DSN from environment or use default
	dsn := getEnv("DATA_SOURCE_NAME", "")
	if dsn == "" {
		// Build from individual components
		host := getEnv("DB_HOST", "localhost")
		dbPort := getEnvInt("DB_PORT", 5432)
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbName := getEnv("DB_NAME", "custapi")
		sslMode := getEnv("DB_SSLMODE", "disable")

		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host, dbPort, user, password, dbName, sslMode,
		)
	}

	return &Config{
		Port:           port,
		DataSourceName: dsn,
		CorsOrigins:    getEnv("CORS_ORIGINS", "http://localhost:3000"),
		Telemetry:      loadTelemetryConfig(),
	}
}

// loadTelemetryConfig reads OTEL settings tuned for shipping traces/metrics
// to an AWS Distro for OpenTelemetry (ADOT) Collector sidecar on ECS, which
// forwards traces to X-Ray and metrics to CloudWatch.
func loadTelemetryConfig() telemetry.Config {
	return telemetry.Config{
		Enabled:              getEnvBool("OTEL_ENABLED", true),
		ServiceName:          getEnv("OTEL_SERVICE_NAME", "custapi"),
		ServiceVersion:       getEnv("OTEL_SERVICE_VERSION", "dev"),
		Environment:          getEnv("DEPLOYMENT_ENVIRONMENT", "development"),
		OTLPEndpoint:         getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		OTLPInsecure:         getEnvBool("OTEL_EXPORTER_OTLP_INSECURE", true),
		TracesSamplerRatio:   getEnvFloat("OTEL_TRACES_SAMPLER_RATIO", 1.0),
		MetricExportInterval: getEnvDuration("OTEL_METRIC_EXPORT_INTERVAL_MS", 60*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

// getEnvDuration reads a millisecond count from the environment (matching
// the OTEL_METRIC_EXPORT_INTERVAL spec convention) and returns it as a
// time.Duration.
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if ms, err := strconv.Atoi(value); err == nil {
			return time.Duration(ms) * time.Millisecond
		}
	}
	return defaultValue
}
