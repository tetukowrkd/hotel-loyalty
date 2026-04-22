package hotel

import (
	stdErrors "errors"

	"database/sql"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/mapper"
	"hotel-loyalty/internal/model/response"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/repository"

	"github.com/google/uuid"
)

type HotelUsecase struct {
	hotelRepo       repository.HotelRepository
	hotelImageRepo  repository.HotelImageRepository
	facilityMapRepo repository.HotelFacilityMapRepository
}

func NewHotelUsecase(hotelRepo repository.HotelRepository, hotelImageRepo repository.HotelImageRepository, facilityMapRepo repository.HotelFacilityMapRepository) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo:       hotelRepo,
		hotelImageRepo:  hotelImageRepo,
		facilityMapRepo: facilityMapRepo,
	}
}

func (u *HotelUsecase) Create(hotel *domain.Hotel) (*response.CreateHotelResponse, error) {

	// 🔥 basic validation (optional, kalau belum di handler)
	if hotel.Name == "" {
		logger.InfoLogger.Println("[HOTEL_USECASE][Create] name is required")
		return nil, errors.ErrBadRequest
	}

	if hotel.StarRating < 1 || hotel.StarRating > 5 {
		logger.InfoLogger.Println("[HOTEL_USECASE][Create] invalid star rating:", hotel.StarRating)
		return nil, errors.ErrBadRequest
	}

	// 🔥 set default value (INI BEDANYA SAMA HANDLER)
	hotel.ID = uuid.New()
	hotel.Status = "draft"

	// 🔥 insert
	err := u.hotelRepo.Create(hotel)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_USECASE][Create] create failed:", err)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[HOTEL_USECASE][Create] success",
		"name=", hotel.Name,
		"city=", hotel.City,
	)

	return mapper.ToHotelResponse(hotel), nil
}

func (u *HotelUsecase) List(city string, star int, limit, offset int) ([]response.HotelListItem, int, error) {

	data, total, err := u.hotelRepo.List(city, star, limit, offset)
	if err != nil {
		logger.ErrorLogger.Println("[USECASE][Hotel][List] failed:", err)
		return nil, 0, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][Hotel][List] success",
		"count", len(data),
		"total", total,
	)

	return data, total, nil
}

func (u *HotelUsecase) GetByID(id uuid.UUID) (*response.HotelDetailResponse, error) {

	h, err := u.hotelRepo.GetByID(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][Hotel][GetByID] hotel not found",
			"hotelID", id,
			"error", err,
		)
		return nil, errors.ErrNotFound
	}

	images, err := u.hotelImageRepo.GetByHotelID(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][Hotel][GetByID] get images failed",
			"hotelID", id,
			"error", err,
		)
		return nil, errors.ErrInternal
	}

	facilities, err := u.facilityMapRepo.GetByHotelID(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][Hotel][GetByID] get facilities failed",
			"hotelID", id,
			"error", err,
		)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][Hotel][GetByID] success",
		"hotelID", id,
		"imageCount", len(images),
		"facilityCount", len(facilities),
	)

	return &response.HotelDetailResponse{
		ID:          h.ID,
		Name:        h.Name,
		Description: h.Description,
		City:        h.City,
		Country:     h.Country,
		StarRating:  h.StarRating,
		Images:      images,
		Facilities:  facilities,
	}, nil
}

func (u *HotelUsecase) Publish(id uuid.UUID) error {

	h, err := u.hotelRepo.GetByID(id)
	if err != nil {
		return errors.ErrNotFound
	}

	if h.Status == "published" {
		logger.InfoLogger.Println(
			"[USECASE][Hotel][Publish] already published",
			"hotelID", id,
		)
		return errors.NewBadRequest("Already published")
	}

	err = u.hotelRepo.Publish(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][Hotel][Publish] failed",
			"hotelID", id,
			"error", err,
		)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][Hotel][Publish] success",
		"hotelID", id,
	)

	return nil
}

func (u *HotelUsecase) Delete(id uuid.UUID) error {

	err := u.hotelRepo.Delete(id)
	if err != nil {
		logger.ErrorLogger.Println("[USECASE][Hotel][Delete] failed:", err)

		if stdErrors.Is(err, sql.ErrNoRows) {
			return errors.ErrNotFound
		}

		return errors.ErrInternal
	}

	logger.InfoLogger.Println("[USECASE][Hotel][Delete] success", "id", id)

	return nil
}
