package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL   string
	ServerAddress string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}

	return cfg
}
