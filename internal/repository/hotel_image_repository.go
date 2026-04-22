package repository

import (
	"hotel-loyalty/internal/domain"

	"github.com/google/uuid"
)

type HotelImageRepository interface {
	Create(image *domain.HotelImage) error
	CountByHotelID(hotelId uuid.UUID) (int, error)
	SetAllNonPrimary(hotelID uuid.UUID) error
	SetPrimary(imageID uuid.UUID) error
	Exists(imageID, hotelID uuid.UUID) (bool, error)
	GetByHotelID(hotelID uuid.UUID) ([]string, error)

	Delete(imageID uuid.UUID) error
}
