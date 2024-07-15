package dto

import "github.com/google/uuid"

type ReqUpdateRole struct {
	Name string `json:"name" validate:"required"`
}

type ToDBUpdateRole struct {
	Name        string `json:"name"`
	UpdatedByID uuid.UUID
}
