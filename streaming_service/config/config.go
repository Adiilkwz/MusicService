package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	AudioDir       string
	NATSUrl        string
	AuthServiceUrl string
	RedisAddr      string
	RedisPassword  string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("STREAMING_SERVER_PORT", "50052"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "music_service"),
		AudioDir:       getEnv("AUDIO_DIR", "./audio_files"),
		NATSUrl:        getEnv("NATS_URL", "nats://localhost:4222"),
		AuthServiceUrl: getEnv("AUTH_SERVICE_URL", "localhost:50051"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		SMTPHost:       getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:       getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
