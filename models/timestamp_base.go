package models

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

type TimestampBase struct {
	CreatedAt   utils.NullTime `json:"created_at"`
	CreatedByID uuid.UUID      `json:"created_by_id"`
	UpdatedAt   utils.NullTime `json:"updated_at"`
	UpdatedByID uuid.UUID      `json:"updated_by_id"`
	DeletedAt   utils.NullTime `json:"deleted_at"`
	DeletedByID uuid.UUID      `json:"deleted_by_id"`
}
