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

	Phone      string `json:"phone"`
	MemberCode string `json:"member_code"`

	LoyaltyPoints int        `json:"loyalty_points"`
	TierID        *uuid.UUID `json:"tier_id"`

	RoleID   *uuid.UUID `json:"role_id"`
	RoleName string     `json:"role_name"`

	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LoginAttempt    int        `json:"login_attempt"`

	ResetPasswordToken     *string    `json:"-"`
	ResetPasswordExpiredAt *time.Time `json:"-"`

	IsActive  bool `json:"is_active"`
	IsDeleted bool `json:"is_deleted"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
