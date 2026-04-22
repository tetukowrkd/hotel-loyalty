package response

import "github.com/google/uuid"

type CreateHotelImageResponse struct {
	ID        uuid.UUID `json:"id"`
	HotelID   uuid.UUID `json:"hotel_id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	SortOrder int       `json:"sort_order"`
}
