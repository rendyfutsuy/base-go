package dto

import "github.com/google/uuid"

type ReqUpdatePermission struct {
	Name string `json:"name" validate:"required"`
}

type ToDBUpdatePermission struct {
	Name        string `json:"name"`
	UpdatedByID uuid.UUID
}
