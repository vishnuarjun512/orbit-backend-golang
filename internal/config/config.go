package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                              string
	AppPort                             string
	DatabaseURL                         string
	JWT_SECRET                          string
	JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES  int
	JWT_REFRESH_TOKEN_EXPIRY_IN_MINUTES int
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Env File not found")
	}

	accessExpiry, err := strconv.Atoi(
		getEnv("JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES", "15"),
	)
	if err != nil {
		log.Fatal("JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES must be an integer")
	}

	refreshExpiry, err := strconv.Atoi(
		getEnv("JWT_REFRESH_TOKEN_EXPIRY_IN_MINUTES", "30"),
	)
	if err != nil {
		log.Fatal("JWT_REFRESH_TOKEN_EXPIRY_IN_MINUTES must be an integer")
	}

	return &Config{
		AppEnv:                              getEnv("APP_ENV", "development"),
		AppPort:                             getEnv("PORT", "8080"),
		DatabaseURL:                         getEnv("DATABASE_URL", ""),
		JWT_SECRET:                          getEnv("JWT_SECRET", ""),
		JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES:  accessExpiry,
		JWT_REFRESH_TOKEN_EXPIRY_IN_MINUTES: refreshExpiry,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
