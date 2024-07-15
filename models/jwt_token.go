package models

import (
	"time"

	"github.com/google/uuid"
)

// User represent the user model
type JWTToken struct {
	UserId      uuid.UUID  `json:"user_id" validate:"required"`
	AccessToken string     `json:"access_token" validate:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
