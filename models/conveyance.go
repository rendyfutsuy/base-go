package models

import "github.com/google/uuid"

type Conveyance struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
	Type string    `json:"type"`
	TimestampBase
}
