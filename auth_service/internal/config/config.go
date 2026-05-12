package config

import (
	"os"
)

type Config struct {
	GRPCPort    string
	DatabaseURL string
	JWTSecret   string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
}

func Load() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/music_service?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		SMTPHost:    getEnv("SMTP_HOST", "localhost"),
		SMTPPort:    getEnv("SMTP_PORT", "2525"),
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
