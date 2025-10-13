package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

// User represent the user model
type User struct {
	ID                uuid.UUID        `json:"id" validate:"required"`
	RoleId            uuid.UUID        `json:"role_id" validate:"required"`
	FullName          string           `json:"full_name" validate:"required"`
	Username          string           `json:"username" validate:"required"`
	Email             string           `json:"email" validate:"required,email"`
	Password          string           `json:"password" validate:"required"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         utils.NullTime   `json:"deleted_at"`
	PasswordExpiredAt time.Time        `json:"password_expired_at"`
	Gender            string           `json:"gender"`
	Counter           int              `json:"counter"`
	ApiKey            utils.NullString `json:"api_key"`

	// mutator
	ActiveStatus utils.NullString `json:"active_status"`
	IsBlocked    bool             `json:"is_blocked"`
	RoleName     string           `json:"role_name"`
}
