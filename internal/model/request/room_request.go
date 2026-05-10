package request

import "github.com/google/uuid"

type CreateRoomRequest struct {
	ID      uuid.UUID `json:"id"`
	HotelID uuid.UUID `json:"hotel_id"`

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
	BasePrice   float64 `json:"base_price"`
}

type SetInventoryRequest struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`
	Stock     int    `json:"stock"`
}
