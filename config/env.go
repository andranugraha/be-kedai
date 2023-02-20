package config

import (
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func GetEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)

	if envValue == "" {
		return defaultValue
	}

	return envValue
}

func GetArrayENV(key string, defaultValue []string) []string {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue
	}

	return strings.Split(env, ",")
}
