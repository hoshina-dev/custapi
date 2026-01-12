package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port           int
	DataSourceName string
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
