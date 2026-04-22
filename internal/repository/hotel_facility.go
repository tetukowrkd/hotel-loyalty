package repository

import (
	"hotel-loyalty/internal/domain"

	"github.com/google/uuid"
)

type HotelFacilityRepository interface {
	GetAll() ([]domain.HotelFacility, error)
	Create(name, icon string) (uuid.UUID, error)
	Delete(id uuid.UUID) error
}

type HotelFacilityMapRepository interface {
	Assign(hotelID uuid.UUID, facilityIDs []uuid.UUID) error
	GetByHotelID(hotelID uuid.UUID) ([]domain.HotelFacility, error)
}
