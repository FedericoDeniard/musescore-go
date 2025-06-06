package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Keys struct {
	ENVIRONMENT             string `json:"enviroment"`
	AWS_DEFAULT_REGION      string `json:"aws_default_region"`
	AWS_USER_POOL_ID        string `json:"aws_user_pool_id"`
	AWS_USER_POOL_CLIENT_ID string `json:"aws_user_pool_client_id"`
}

func getKeys() Keys {
	if os.Getenv("ENVIROMENT") == "production" {
		return Keys{
			ENVIRONMENT:             os.Getenv("ENVIROMENT"),
			AWS_DEFAULT_REGION:      os.Getenv("AWS_DEFAULT_REGION"),
			AWS_USER_POOL_ID:        os.Getenv("AWS_USER_POOL_ID"),
			AWS_USER_POOL_CLIENT_ID: os.Getenv("AWS_USER_POOL_CLIENT_ID"),
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Keys{
		ENVIRONMENT:             os.Getenv("ENVIROMENT"),
		AWS_DEFAULT_REGION:      os.Getenv("AWS_DEFAULT_REGION"),
		AWS_USER_POOL_ID:        os.Getenv("AWS_USER_POOL_ID"),
		AWS_USER_POOL_CLIENT_ID: os.Getenv("AWS_USER_POOL_CLIENT_ID"),
	}
}

var KEYS = getKeys()
