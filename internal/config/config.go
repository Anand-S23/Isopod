package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// MinSessionKeyLength is the minimum size for gorilla/secure cookie signing / encryption.
const MinSessionKeyLength = 32

type EnvVars struct {
	PRODUCTION           bool
	PORT                 string
	DB_URI               string
	ALLOWED_ORIGIN       string
	ENCRYPTION_KEY       []byte
	SESSION_KEY          []byte
	CSRF_TOKEN           string
	BASE_URL             string
	GITHUB_CLIENT_ID     string
	GITHUB_CLIENT_SECRET string
}

func LoadEnv() *EnvVars {
	envMode := GetEnv("MODE", "development")

	sessionKey := []byte(GetEnvOrPanic("SESSION_KEY"))
	if len(sessionKey) < MinSessionKeyLength {
		log.Fatalf("SESSION_KEY must be at least %d bytes (got %d)", MinSessionKeyLength, len(sessionKey))
	}

	return &EnvVars{
		PRODUCTION:           envMode == "production",
		PORT:                 GetEnv("PORT", "8080"),
		DB_URI:               GetEnvOrPanic("DB_URI"),
		ALLOWED_ORIGIN:       GetEnvOrPanic("ALLOWED_ORIGIN"),
		ENCRYPTION_KEY:       []byte(GetEnvOrPanic("ENCRYPTION_KEY")),
		SESSION_KEY:          sessionKey,
		CSRF_TOKEN:           GetEnvOrPanic("CSRF_TOKEN"),
		BASE_URL:             GetEnvOrPanic("BASE_URL"),
		GITHUB_CLIENT_ID:     GetEnvOrPanic("GITHUB_CLIENT_ID"),
		GITHUB_CLIENT_SECRET: GetEnvOrPanic("GITHUB_CLIENT_SECRET"),
	}
}

func GetEnv(env string, defaultValue string) string {
	variable := os.Getenv(env)
	if variable == "" {
		return strings.TrimSpace(defaultValue)
	}
	return strings.TrimSpace(variable)
}

func GetEnvOrPanic(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		message := fmt.Sprintf("must provide %s variable in .env file", env)
		log.Fatal(message)
	}
	return strings.TrimSpace(variable)
}
