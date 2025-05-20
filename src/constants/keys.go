package constants

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENVIROMENT string
}

var KEYS Config

func init() {
	godotenv.Load()
	ENVIRONMENT := os.Getenv("ENVIROMENT")

	if ENVIRONMENT == "" {
		ENVIRONMENT = "development"
	}

	KEYS = Config{
		ENVIROMENT: ENVIRONMENT,
	}

}
