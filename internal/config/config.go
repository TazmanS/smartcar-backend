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

	MQTTHost          string
	MQTTPort          string
	MQTTCarSessionKey string

	ESP32URL    string
	BackendHost string
	BackendPort string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		PORT:              getEnv("PORT", ":8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBName:            getEnv("DB_NAME", "smartcar"),
		DBUser:            getEnv("DB_USER", "admin"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		MQTTHost:          getEnv("MQTT_HOST", "localhost"),
		MQTTPort:          getEnv("MQTT_PORT", "1883"),
		MQTTCarSessionKey: getEnv("MQTT_CAR_SESSION_KEY", "smartcar-esp32"),
		ESP32URL:          getEnv("ESP32_URL", "http://192.168.31.111"),
		BackendHost:       getEnv("BACKEND_HOST", "http://192.168.31.210"),
		BackendPort:       getEnv("BACKEND_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
