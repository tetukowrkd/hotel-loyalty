package repository

import "hotel-loyalty/internal/domain"

type RoleRepository interface {
	GetByName(name string) (*domain.Role, error)
}
