package models

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	TimestampBase
}
