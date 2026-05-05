package domain

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID      uuid.UUID `json:"id"`
	HotelID uuid.UUID `json:"hotel_id"`

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	BasePrice   float64 `json:"base_price"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
