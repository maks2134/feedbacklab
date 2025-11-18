package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	HealthPort    string
	DatabaseURL   string
	MigrationsDir string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:        getEnv("APP_PORT", "8080"),
		HealthPort:     getEnv("HEALTH_PORT", "8081"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://feedback:feedback@db:5432/innotech?sslmode=disable"),
		MigrationsDir:  getEnv("MIGRATIONS_DIR", "./migrations"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:    getEnv("MINIO_BUCKET", "feedbacklab"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
	}

	log.Println("config loaded successfully")
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
