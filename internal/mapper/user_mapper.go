package mapper

import (
	"hotel-loyalty/internal/domain"
	"hotel-loyalty/internal/model/request"
	"hotel-loyalty/internal/model/response"
)

func ToUserDomain(req request.RegisterUserRequest) *domain.User {
	return &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		RoleName: req.RoleName,
	}
}

func ToRegisterResponse(u *domain.User) *response.RegisterUserResponse {
	return &response.RegisterUserResponse{
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		MemberCode: u.MemberCode,
	}
}

func ToUserResponse(u *domain.User) *response.UserResponse {
	return &response.UserResponse{
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		Role:       u.RoleName,
		MemberCode: u.MemberCode,
	}
}
