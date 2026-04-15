package repository

import "hotel-loyalty/internal/domain"

type TokenRepository interface {
	Save(token *domain.UserToken) error
	FindByRefreshToken(token string) (*domain.UserToken, error)
	DeleteByRefreshToken(token string) error
}
