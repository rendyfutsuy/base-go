package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

// Permission represent the role model
type Permission struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
	DeletedAt utils.NullTime `json:"deleted_at"`

	// relation
	PermissionGroups []PermissionGroup `json:"permission_groups"`
	Roles            []Role            `json:"roles"`
}
