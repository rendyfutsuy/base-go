package dto

import "github.com/google/uuid"

type ToDBDeleteAccount struct {
	DeletedByID uuid.UUID `json:"deleted_by"`
}
