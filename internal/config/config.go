package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type EnvVars struct {
	PRODUCTION   bool
	PORT         string
	DATABASE_URL string
}

func LoadEnv() (*EnvVars, error) {
	envMode := GetEnv("MODE", "development")
	port := GetEnv("PORT", "8080")
	databaseURL := GetEnvOrPanic("DB_URI")

	return &EnvVars{
		PRODUCTION:   envMode == "production",
		PORT:         port,
		DATABASE_URL: databaseURL,
	}, nil
}

func GetEnv(env string, defaultValue string) string {
	variable := os.Getenv(env)
	if variable == "" {
		return defaultValue
	}
	return variable
}

func GetEnvOrPanic(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		message := fmt.Sprintf("must provide %s variable in .env file", env)
		log.Fatal(message)
	}
	return strings.TrimSpace(variable)
}
