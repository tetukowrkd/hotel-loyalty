package main

import (
	"log"

	"github.com/joho/godotenv"

	"ecommerce-be/internal/config"
	"ecommerce-be/internal/infrastructure/database"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// load config
	cfg := config.LoadConfig()

	// connect DB
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}