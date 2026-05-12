package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GRPCPort string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string

	NatsURL string
}

func Load() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "music_catalog"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		NatsURL: getEnv("NATS_URL", "nats://localhost:4222"),
	}
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
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
