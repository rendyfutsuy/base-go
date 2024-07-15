package models

import (
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

// Shipyard represents a shipyard with its details.
type Shipyard struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Yard      string         `json:"yard"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt utils.NullTime `json:"deleted_at"`
}
