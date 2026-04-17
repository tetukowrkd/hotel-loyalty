package main

import (
	"github.com/joho/godotenv"

	"hotel-loyalty/internal/app"
	"hotel-loyalty/internal/config"
	"hotel-loyalty/internal/infrastructure/jwt"
	"hotel-loyalty/internal/infrastructure/logger"
)

func main() {
	// init logger
	logger.InitLogger()

	logger.InfoLogger.Println("[MAIN] starting application")

	// load env (dev only)
	if err := godotenv.Load(); err != nil {
		logger.InfoLogger.Println("[MAIN] .env not found, using system env")
	}

	// load config
	cfg := config.LoadConfig()

	// validate config
	if err := cfg.Validate(); err != nil {
		logger.ErrorLogger.Println("[MAIN] config error:", err)
		panic(err)
	}

	// init JWT
	jwt.InitJWT(cfg.JWTSecret, cfg.JWTExp)
	logger.InfoLogger.Println("[MAIN] JWT initialized")

	// init container
	container := app.NewContainer(cfg)

	// start app
	application := app.NewApp(container)
	application.Start()
}
