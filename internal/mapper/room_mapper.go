package mapper

import (
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
)

func ToRoomDomain(req request.CreateRoomRequest) *domain.Room {
	return &domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
		BasePrice:   req.BasePrice,
	}
}

func ToRoomListItem(u domain.Room) response.RoomListItem {
	return response.RoomListItem{
		ID:        u.ID,
		Name:      u.Name,
		Capacity:  u.Capacity,
		BasePrice: u.BasePrice,
	}
}

func ToRoomResponse(u *domain.Room) *response.CreateRoomResponse {
	return &response.CreateRoomResponse{
		ID:          u.ID,
		HotelID:     u.HotelID,
		Name:        u.Name,
		Description: u.Description,
		Capacity:    u.Capacity,
		BasePrice:   u.BasePrice,
	}
}
