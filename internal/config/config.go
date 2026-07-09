package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		PORT: getEnv("PORT", ":8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
