package response

import (
	"hotel-loyalty/internal/domain"

	"github.com/google/uuid"
)

type CreateHotelResponse struct {
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

type HotelListItem struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	City       string    `json:"city"`
	Country    string    `json:"country"`
	StarRating int       `json:"star_rating"`
	ImageURL   string    `json:"image_url"`
}

type HotelDetailResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	City        string                 `json:"city"`
	Country     string                 `json:"country"`
	StarRating  int                    `json:"star_rating"`
	Images      []string               `json:"images"`
	Facilities  []domain.HotelFacility `json:"facilities"`
}
