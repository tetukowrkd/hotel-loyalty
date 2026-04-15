package domain

import (
	"time"
)

type UserToken struct {
	ID           string
	UserID       string
	RefreshToken string
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
}
