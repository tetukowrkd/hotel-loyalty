package request

import "github.com/google/uuid"

type CreateHotelRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Lat         float64   `json:"lat"`
	Lng         float64   `json:"lng"`
	StarRating  int       `json:"star_rating"`
	Status      string    `json:"status"`
}
