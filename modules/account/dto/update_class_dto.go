package dto

import "github.com/google/uuid"

type ReqUpdateAccount struct {
	Name string `json:"name" validate:"required"`
}

type ToDBUpdateAccount struct {
	Name        string `json:"name"`
	UpdatedByID uuid.UUID
}
