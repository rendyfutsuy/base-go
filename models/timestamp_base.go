package models

import "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"

type TimestampBase struct {
	CreatedAt   utils.NullTime   `json:"created_at"`
	CreatedByID utils.NullString `json:"created_by_id"`
	UpdatedAt   utils.NullTime   `json:"updated_at"`
	UpdatedByID utils.NullString `json:"updated_by_id"`
	DeletedAt   utils.NullTime   `json:"deleted_at"`
	DeletedByID utils.NullString `json:"deleted_by_id"`
}
