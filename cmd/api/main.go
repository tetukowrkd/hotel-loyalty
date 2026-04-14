package main

import (
	"github.com/joho/godotenv"

	"e-commerce/internal/app"
	"e-commerce/internal/config"
	"e-commerce/internal/infrastructure/logger"
)

func main() {
	// init logger
	logger.InitLogger()

	// load env
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLogger.Println("No .env file found")
	}

	// load config
	cfg := config.LoadConfig()

	// start app
	application := app.NewApp(cfg)
	application.Start()
}
