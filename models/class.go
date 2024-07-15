package models

import "github.com/google/uuid"

type Class struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
	TimestampBase
}
