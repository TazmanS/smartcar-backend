package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	FrontendURL string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}
}
