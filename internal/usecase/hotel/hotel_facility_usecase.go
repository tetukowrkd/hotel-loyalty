package hotel

import (
	"database/sql"
	stdErrors "errors"
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/infrastructure/logger"
	"hotel-loyalty/internal/pkg/errors"
	"hotel-loyalty/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type HotelFacilityUsecase struct {
	facilityRepo repository.HotelFacilityRepository
	mapRepo      repository.HotelFacilityMapRepository
}

func NewHotelFacilityUsecase(facilityRepo repository.HotelFacilityRepository, mapRepo repository.HotelFacilityMapRepository) *HotelFacilityUsecase {
	return &HotelFacilityUsecase{
		facilityRepo: facilityRepo,
		mapRepo:      mapRepo,
	}
}

func (u *HotelFacilityUsecase) Assign(hotelID uuid.UUID, facilityIDs []uuid.UUID) error {

	if len(facilityIDs) == 0 {
		return errors.NewBadRequest("Facilities cannot be empty")
	}

	err := u.mapRepo.Assign(hotelID, facilityIDs)
	if err != nil {
		logger.ErrorLogger.Println("[HOTEL_FACILITY_USECASE] assign failed:", err)
		return errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[HOTEL_FACILITY_USECASE] assign success",
		"hotelID", hotelID,
		"count", len(facilityIDs),
	)

	return nil
}

func (u *HotelFacilityUsecase) GetByHotelID(hotelID uuid.UUID) ([]domain.HotelFacility, error) {

	data, err := u.mapRepo.GetByHotelID(hotelID)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][HotelFacility][GetByHotelID] failed",
			"hotelID", hotelID,
			"error", err,
		)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][HotelFacility][GetByHotelID] success",
		"hotelID", hotelID,
		"count", len(data),
	)

	return data, nil
}

func (u *HotelFacilityUsecase) GetAll() ([]domain.HotelFacility, error) {

	data, err := u.facilityRepo.GetAll()
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][HotelFacility][GetAll] failed",
			"error", err,
		)
		return nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][HotelFacility][GetAll] success",
		"count", len(data),
	)

	return data, nil
}

func (u *HotelFacilityUsecase) Create(name, icon string) (uuid.UUID, error) {

	name = strings.TrimSpace(name)
	if name == "" {
		return uuid.Nil, errors.NewBadRequest("Name is required")
	}

	if icon == "" {
		return uuid.Nil, errors.NewBadRequest("Icon is required")
	}

	id, err := u.facilityRepo.Create(name, icon)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][HotelFacility][Create] failed",
			"name", name,
			"error", err,
		)
		return uuid.Nil, errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][HotelFacility][Create] success",
		"id", id,
		"name", name,
	)

	return id, nil
}

func (u *HotelFacilityUsecase) Delete(id uuid.UUID) error {

	err := u.facilityRepo.Delete(id)
	if err != nil {
		logger.ErrorLogger.Println(
			"[USECASE][HotelFacility][Delete] failed",
			"facilityID", id,
			"error", err,
		)

		if stdErrors.Is(err, sql.ErrNoRows) {
			return errors.ErrNotFound
		}

		return errors.ErrInternal
	}

	logger.InfoLogger.Println(
		"[USECASE][HotelFacility][Delete] success",
		"facilityID", id,
	)

	return nil
}
