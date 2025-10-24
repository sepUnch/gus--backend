package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_DSN     string
	JWT_SECRET string
)

func LoadConfig() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Could not load .env file, using environment variables")
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")

	JWT_SECRET = os.Getenv("JWT_SECRET")

	DB_DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	if JWT_SECRET == "" {
		log.Fatal("FATAL: JWT_SECRET not set in .env file or environment")
	}
	if dbUser == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("FATAL: Database credentials (DB_HOST, DB_PORT, DB_USERNAME, DB_DATABASE) not set completely in .env file or environment")
	}
}
