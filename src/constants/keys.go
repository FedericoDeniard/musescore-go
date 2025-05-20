package constants

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENVIRONMENT string
}

var KEYS Config

func init() {
	godotenv.Load()
	environment := os.Getenv("ENVIRONMENT")

	if environment == "" {
		environment = "development"
	}

	KEYS = Config{
		ENVIRONMENT: environment,
	}

}
