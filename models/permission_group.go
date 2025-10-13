package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

type PermissionGroup struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Module    utils.NullString `json:"module"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt utils.NullTime   `json:"updated_at"`
	DeletedAt utils.NullTime   `json:"deleted_at"`

	// relation
	Permissions     []Permission       `json:"permissions"`
	PermissionNames []utils.NullString `json:"permission_names"`
}
