package models

import (
	"time"

	"github.com/google/uuid"
)

// User represent the user model
type PasswordHistory struct {
	UserId         uuid.UUID  `json:"user_id" validate:"required"`
	HashedPassword string     `json:"hashed_password" validate:"required"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
