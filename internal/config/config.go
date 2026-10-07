package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	DatabaseURL string
	JWT_SECRET  string
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Env File not found")
	}

	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		AppPort:     getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWT_SECRET:  getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
