// Package config provides configuration loading and management for the application.
package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration settings.
type Config struct {
	AppPort       int
	HealthPort    int
	DatabaseURL   string
	ParsedDBURL   *url.URL
	MigrationsDir string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

// Load reads configuration from environment variables and returns a Config instance.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	appPort, err := getEnvInt("APP_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %w", err)
	}
	cfg.AppPort = appPort

	healthPort, err := getEnvInt("HEALTH_PORT", 8081)
	if err != nil {
		return nil, fmt.Errorf("invalid HEALTH_PORT: %w", err)
	}
	cfg.HealthPort = healthPort

	cfg.DatabaseURL = getEnv("DATABASE_URL",
		"postgres://feedback:feedback@db:5432/innotech?sslmode=disable",
	)

	parsedURL, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL: %w", err)
	}
	cfg.ParsedDBURL = parsedURL

	cfg.MigrationsDir = getEnv("MIGRATIONS_DIR", "./migrations")

	cfg.MinioEndpoint = getEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.MinioAccessKey = getEnv("MINIO_ACCESS_KEY", "minio")
	cfg.MinioSecretKey = getEnv("MINIO_SECRET_KEY", "minio123")
	cfg.MinioBucket = getEnv("MINIO_BUCKET", "feedback")

	useSSLStr := getEnv("MINIO_USE_SSL", "false")
	cfg.MinioUseSSL, _ = strconv.ParseBool(useSSLStr)

	log.Println("config loaded and parsed successfully")
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

// getEnvInt - parser to integer value
func getEnvInt(key string, fallback int) (int, error) {
	if val, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}
	return fallback, nil
}
