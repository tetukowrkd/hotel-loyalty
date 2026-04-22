package app

import (
	"hotel-loyalty/internal/config"
	handlerHotel "hotel-loyalty/internal/delivery/http/handler/hotel"
	handlerIdentity "hotel-loyalty/internal/delivery/http/handler/identity"
	"hotel-loyalty/internal/infrastructure/database"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/repository/postgres"
	usecaseHotel "hotel-loyalty/internal/usecase/hotel"
	usecaseIdentity "hotel-loyalty/internal/usecase/identity"
)

type Container struct {
	UserUsecase          *usecaseIdentity.UserUsecase
	UserHandler          *handlerIdentity.UserHandler
	TokenUsecase         *usecaseIdentity.TokenUsecase
	TokenHandler         *handlerIdentity.TokenHandler
	HotelUsecase         *usecaseHotel.HotelUsecase
	HotelHandler         *handlerHotel.HotelHandler
	HotelFacilityUsecase *usecaseHotel.HotelFacilityUsecase
	HotelFacilityHandler *handlerHotel.HotelFacilityHandler
	HotelImageUsecase    *usecaseHotel.HotelImageUsecase
	HotelImageHandler    *handlerHotel.HotelImageHandler
}

func NewContainer(cfg *config.Config) *Container {
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.ErrorLogger.Fatal("[CONTAINER] DB error:", err)
	}

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	hotelRepo := postgres.NewHotelRepository(db)
	hotelImageRepo := postgres.NewHotelImageRepository(db)
	hotelFacilityRepo := postgres.NewHotelFacilityRepository(db)
	hotelFacilityMapRepo := postgres.NewHotelFacilityMapRepository(db)

	// Auth + User
	userUsecase := usecaseIdentity.NewUserUsecase(
		userRepo,
		tokenRepo,
		roleRepo,
		cfg.RefreshTokenExp,
	)
	tokenUsecase := usecaseIdentity.NewTokenUsecase(
		userRepo,
		tokenRepo,
		cfg.RefreshTokenExp,
	)

	// Hotel
	hotelUsecase := usecaseHotel.NewHotelUsecase(
		hotelRepo,
		hotelImageRepo,
		hotelFacilityMapRepo,
	)
	hotelImageUsecase := usecaseHotel.NewHotelImageUsecase(hotelImageRepo, db)
	hotelFacilityUsecase := usecaseHotel.NewHotelFacilityUsecase(
		hotelFacilityRepo,
		hotelFacilityMapRepo,
	)

	userHandler := handlerIdentity.NewUserHandler(userUsecase)
	tokenHandler := handlerIdentity.NewTokenHandler(tokenUsecase)
	hotelHandler := handlerHotel.NewHotelHandler(hotelUsecase)
	hotelImageHandler := handlerHotel.NewHotelImageHandler(hotelImageUsecase)
	hotelFacilityHandler := handlerHotel.NewHotelFacilityHandler(hotelFacilityUsecase)

	return &Container{
		UserHandler:          userHandler,
		TokenHandler:         tokenHandler,
		HotelHandler:         hotelHandler,
		HotelImageHandler:    hotelImageHandler,
		HotelFacilityHandler: hotelFacilityHandler,
	}
}
