package mapper

import (
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/model/response"
	"os"
)

func ToHotelImageResponse(u *domain.HotelImage) *response.CreateHotelImageResponse {
	baseURL := os.Getenv("BASE_URL")

	imageURL := u.ImageURL
	if imageURL != "" {
		imageURL = baseURL + imageURL
	}

	return &response.CreateHotelImageResponse{
		ID:        u.ID,
		HotelID:   u.HotelID,
		ImageURL:  imageURL,
		IsPrimary: u.IsPrimary,
		SortOrder: u.SortOrder,
	}
}
