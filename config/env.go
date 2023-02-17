package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)

	if envValue == "" {
		return defaultValue
	}

	return envValue
}
