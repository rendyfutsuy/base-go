package models

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

// Permission represent the role model
type Permission struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Deletable bool      `json:"deletable"`
	TotalUser *int      `json:"total_user"`

	Users     []User         `json:"users"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}
