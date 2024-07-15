package dto

import "github.com/google/uuid"

type ToDBDeleteRole struct {
	DeletedByID uuid.UUID `json:"deleted_by"`
}
