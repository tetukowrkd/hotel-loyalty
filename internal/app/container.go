package app

import (
	"hotel-loyalty/internal/config"
	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/infrastructure/database"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/repository/postgres"
	"hotel-loyalty/internal/usecase"
)

type Container struct {
	UserUsecase *usecase.UserUsecase
	UserHandler *handler.UserHandler
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatal("[CONTAINER] DB error:", err)
	}

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	userUsecase := usecase.NewUserUsecase(
		userRepo,
		tokenRepo,
		cfg.RefreshTokenExp,
	)

	userHandler := handler.NewUserHandler(userUsecase)

	return &Container{
		UserHandler: userHandler,
	}
}
