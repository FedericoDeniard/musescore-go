package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Keys struct {
	ENVIRONMENT string `json:"api_token"`
}

func getKeys() Keys {
	if os.Getenv("ENVIROMENT") == "production" {
		return Keys{
			ENVIRONMENT: os.Getenv("ENVIROMENT"),
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Keys{
		ENVIRONMENT: os.Getenv("ENVIROMENT"),
	}
}

var KEYS = getKeys()
