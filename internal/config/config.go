package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() Config {
	return Config{
		DatabaseURL: getEnv("DB_URL", ""),
		Port:        getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
