package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

// Permission represent the permission model
type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Deletable bool           `gorm:"column:deletable;default:true;not null" json:"deletable"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt utils.NullTime `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt utils.NullTime `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relations
	PermissionGroups []PermissionGroup `gorm:"many2many:permissions_modules;" json:"permission_groups"`
	Roles            []Role            `gorm:"many2many:permissions_modules;" json:"roles"`
}

// TableName specifies table name for GORM
func (Permission) TableName() string {
	return "permissions"
}
