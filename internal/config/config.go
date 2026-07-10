package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		PORT:       getEnv("PORT", ":8080"),
		DBHost:     getEnv("DBHost", "localhost"),
		DBPort:     getEnv("DBPort", "5432"),
		DBName:     getEnv("DBName", "smartcar"),
		DBUser:     getEnv("DBUser", "admin"),
		DBPassword: getEnv("DBPassword", "password"),
		DBSSLMode:  getEnv("DBSSLMode", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
