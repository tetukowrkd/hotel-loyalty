package app

import (
	"database/sql"

	"hotel-loyalty/internal/delivery/http/handler"
	"hotel-loyalty/internal/repository/postgres"
	"hotel-loyalty/internal/usecase"
)

type Container struct {
	UserHandler *handler.UserHandler
}

func NewContainer(db *sql.DB) *Container {
	// repo
	userRepo := postgres.NewUserRepository(db)

	// usecase
	userUsecase := usecase.NewUserUsecase(userRepo)

	// handler
	userHandler := handler.NewUserHandler(userUsecase)

	return &Container{
		UserHandler: userHandler,
	}
}
