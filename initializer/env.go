package initializer

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var RequiredEnvVars []string = []string{"APP_ENV", "DB_USER", "DB_PASSWORD", "DB_DOMAIN", "DB_NAME", "CLIENT_SECRET", "LOG_LEVEL", "BASE_URL"}

func LoadEnvVariables() {
	currentEnv := os.Getenv("APP_ENV")
	log.Println("Current environment: " + currentEnv)

	/*
		If currentEnv exists, that means this project currently runs in docker.
		Therefore, only load env when the environment is development.

		If it doesn't, that means this project doesn't run in docker.
		Still load env regardless what the current environment is

		TODO: improve this
	*/
	shouldCheckEnv := currentEnv == "" || currentEnv == "development"

	if shouldCheckEnv {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func CheckAllEnvironmentVariables() {
	for _, envKey := range RequiredEnvVars {
		if os.Getenv(envKey) == "" {
			log.Fatalf("[Error] mandatory key %s is missing", envKey)
		}
	}
}
