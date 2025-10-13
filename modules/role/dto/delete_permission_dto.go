package dto

import "github.com/google/uuid"

type ToDBDeletePermission struct {
	DeletedByID uuid.UUID `json:"deleted_by"`
}
