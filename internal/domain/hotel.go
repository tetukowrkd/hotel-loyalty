package domain

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
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

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
