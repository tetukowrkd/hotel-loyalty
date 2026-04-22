package mapper

import (
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
)

func ToHotelDomain(req request.CreateHotelRequest) *domain.Hotel {
	return &domain.Hotel{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		Country:     req.Country,
		Lat:         req.Lat,
		Lng:         req.Lng,
		StarRating:  req.StarRating,
	}
}

func ToHotelResponse(u *domain.Hotel) *response.CreateHotelResponse {
	return &response.CreateHotelResponse{
		ID:          u.ID,
		Name:        u.Name,
		Description: u.Description,
		Address:     u.Address,
		City:        u.City,
		Country:     u.Country,
		Lat:         u.Lat,
		Lng:         u.Lng,
		StarRating:  u.StarRating,
		Status:      u.Status,
	}
}
