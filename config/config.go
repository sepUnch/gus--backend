package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_DSN     string
	JWT_SECRET string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using environment variables")
	}

	DB_DSN = os.Getenv("DB_DSN")
	JWT_SECRET = os.Getenv("JWT_SECRET")

	if DB_DSN == "" {
		log.Fatal("DB_DSN not set")
	}
	if JWT_SECRET == "" {
		log.Fatal("JWT_SECRET not set")
	}
}
