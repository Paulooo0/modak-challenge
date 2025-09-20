package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func Load() Config {
	godotenv.Load()
	return Config{
		DatabaseURL: getEnv("DB_URL", ""),
		Port:        getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Println("Environment variable " + key + " not found, using fallback value: " + fallback)
	return fallback
}
