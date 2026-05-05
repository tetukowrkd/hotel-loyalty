package response

import (
	"github.com/google/uuid"
)

type CreateRoomResponse struct {
	ID      uuid.UUID `json:"id"`
	HotelID uuid.UUID `json:"hotel_id"`

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	BasePrice   float64 `json:"base_price"`
}

type RoomListItem struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	BasePrice float64   `json:"base_price"`
}
