package models

import "github.com/google/uuid"

type Occupation struct {
	ID          uuid.UUID `json:"id"`
	SubCobID    uuid.UUID `json:"subcob_id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	TimestampBase
}
