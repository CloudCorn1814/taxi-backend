package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

type Config struct {
Port string	
DatabaseURL string
}

func LoadConfig() *Config {
if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if port[0] != ':'{
		port = ":" + port
	}

    databaseURL := os.Getenv("DATABASE_URL")

	return &Config {
		Port: port,
		DatabaseURL: databaseURL,
	}
}
	


