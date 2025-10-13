package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

// Role represent the role model
type Role struct {
	ID                   uuid.UUID          `json:"id"`
	Name                 string             `json:"name" validate:"required"`
	Deletable            bool               `json:"deletable"`
	TotalUser            int                `json:"total_user"`
	Modules              []utils.NullString `json:"modules"`
	CategoryNames        []utils.NullString `json:"category_names"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            utils.NullTime     `json:"updated_at"`
	DeletedAt            utils.NullTime     `json:"deleted_at"`
	PermissionGroupNames []utils.NullString `json:"permission_group_names"`
	PermissionGroupIds   []uuid.UUID        `json:"permission_group_ids"`
	Description          utils.NullString   `json:"description"`

	// relation
	Users            []User            `json:"users"`
	Permissions      []Permission      `json:"permissions"`
	PermissionGroups []PermissionGroup `json:"permission_groups"`
}
