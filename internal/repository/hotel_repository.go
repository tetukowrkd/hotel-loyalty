package repository

import (
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/model/response"

	"github.com/google/uuid"
)

type HotelRepository interface {
	Create(hotel *domain.Hotel) error

	List(city string, star int, limit, offset int) ([]response.HotelListItem, int, error)
	GetByID(id uuid.UUID) (*domain.Hotel, error)

	Publish(id uuid.UUID) error

	Delete(id uuid.UUID) error
}
