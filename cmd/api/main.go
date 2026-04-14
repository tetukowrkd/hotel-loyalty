package main

import (
	"github.com/joho/godotenv"

	"e-commerce/internal/config"
	"e-commerce/internal/infrastructure/database"
	"e-commerce/internal/infrastructure/logger"
)

func main() {
	// init logger sekali saja
	logger.InitLogger()

	// load env
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLogger.Println("No .env file found")
	}

	// load config
	cfg := config.LoadConfig()

	logger.InfoLogger.Println("Starting application...")

	// connect DB
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatal("DB connection failed:", err)
	}

	defer db.Close()

	logger.InfoLogger.Println("Application started successfully")
}
