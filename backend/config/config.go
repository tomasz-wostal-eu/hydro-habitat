package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	APIPort     string
	GinMode     string
	JWTSecret   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/hydro_habitat?sslmode=disable"),
		APIPort:     getEnv("API_PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		JWTSecret:   getEnv("JWT_SECRET", "defaultsecret"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
