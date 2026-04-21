package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type EnvVars struct {
	PRODUCTION     bool
	PORT           string
	DB_URI         string
	ALLOWED_ORIGIN string
}

func LoadEnv() (*EnvVars, error) {
	envMode := GetEnv("MODE", "development")

	return &EnvVars{
		PRODUCTION:     envMode == "production",
		PORT:           GetEnv("PORT", "8080"),
		DB_URI:         GetEnvOrPanic("DB_URI"),
		ALLOWED_ORIGIN: GetEnvOrPanic("ALLOWED_ORIGIN"),
	}, nil
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
