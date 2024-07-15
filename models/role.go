package models

import "github.com/google/uuid"

// Role represent the role model
type Role struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Deletable bool      `json:"deletable"`
	TotalUser *int      `json:"total_user"`

	Users       []User       `json:"users"`

}
