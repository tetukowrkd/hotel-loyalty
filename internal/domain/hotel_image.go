package domain

import "github.com/google/uuid"

type HotelImage struct {
	ID        uuid.UUID `json:"id"`
	HotelID   uuid.UUID `json:"hotel_id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
}
