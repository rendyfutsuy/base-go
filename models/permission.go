package models

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/utils"
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
