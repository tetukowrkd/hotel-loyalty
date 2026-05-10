package repository

import (
	"hotel-loyalty/internal/domain"
	"time"

	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(room *domain.Room) error
	ListByHotelID(hotelID uuid.UUID) ([]domain.Room, error)
	Delete(roomID uuid.UUID) error
	UpsertInventory(roomID uuid.UUID, date time.Time, stock int) error
	Exists(roomID uuid.UUID) (bool, error)
}
