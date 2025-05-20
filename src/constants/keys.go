package constants

import (
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	Port string
}

var KEYS Config

func init(){
	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	KEYS = Config{
		Port: port,
	}

}