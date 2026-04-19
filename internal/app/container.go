package app

import (
	"hotel-loyalty/internal/config"
	handlerIdentity "hotel-loyalty/internal/delivery/http/handler/identity"
	"hotel-loyalty/internal/infrastructure/database"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/repository/postgres"
	usecaseIdentity "hotel-loyalty/internal/usecase/identity"
)

type Container struct {
	UserUsecase  *usecaseIdentity.UserUsecase
	UserHandler  *handlerIdentity.UserHandler
	TokenUsecase *usecaseIdentity.TokenUsecase
	TokenHandler *handlerIdentity.TokenHandler
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatal("[CONTAINER] DB error:", err)
	}

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	userUsecase := usecaseIdentity.NewUserUsecase(
		userRepo,
		tokenRepo,
		cfg.RefreshTokenExp,
	)

	tokenUsecase := usecaseIdentity.NewTokenUsecase(
		userRepo,
		tokenRepo,
		cfg.RefreshTokenExp,
	)

	userHandler := handlerIdentity.NewUserHandler(userUsecase)
	tokenHandler := handlerIdentity.NewTokenHandler(tokenUsecase)

	return &Container{
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
	}
}
