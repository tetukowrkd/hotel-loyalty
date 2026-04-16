package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"` // 🔥 jangan expose

	RoleID   *uuid.UUID `json:"role_id"`
	RoleName string     `json:"role_name"`

	Phone string `json:"phone"`

	IsActive  bool `json:"is_active"`
	IsDeleted bool `json:"is_deleted"`

	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`

	LoginAttempt int `json:"login_attempt"`

	ResetPasswordToken     *string    `json:"-"`
	ResetPasswordExpiredAt *time.Time `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
