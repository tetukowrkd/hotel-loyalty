package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
}
