package config

import (
	"fmt"
	"log"
	"os"
)

type EnvVars struct {
    PRODUCTION bool
    PORT       string
}

func LoadEnv() (*EnvVars, error) {
    envMode := GetEnv("MODE", "development")
    port    := GetEnv("PORT", "8080")

    return &EnvVars {
        PRODUCTION: (envMode == "production"),
        PORT: port,
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
	return variable
} 

