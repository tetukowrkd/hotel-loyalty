package repository

import (
	"hotel-loyalty/internal/domain"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(room *domain.Room) error
	ListByHotelID(hotelID uuid.UUID) ([]domain.Room, error)
	Delete(roomID uuid.UUID) error
}
