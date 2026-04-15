package main

import (
	"github.com/joho/godotenv"

	"hotel-loyalty/internal/app"
	"hotel-loyalty/internal/config"
	"hotel-loyalty/internal/infrastructure/logger"
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
