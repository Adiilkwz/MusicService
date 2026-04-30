package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	AudioDir      string
}

func Load() *Config {
	return &Config{
		ServerPort:    getEnv("STREAMING_SERVER_PORT", "50051"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnvAsInt("DB_PORT", 5432),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "music_service"),
		AudioDir:      getEnv("AUDIO_DIR", "./audio_files"),
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