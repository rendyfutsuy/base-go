package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

type PermissionGroup struct {
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	Name        string           `gorm:"type:varchar(255);not null" json:"name"`
	Module      utils.NullString `gorm:"column:module;type:varchar(255);not null" json:"module"`
	Description utils.NullString `gorm:"column:description;type:text" json:"description"`
	Deletable   bool             `gorm:"column:deletable;default:true;not null" json:"deletable"`
	CreatedAt   time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   utils.NullTime   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   utils.NullTime   `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Computed/virtual field - not stored in DB
	PermissionNames []utils.NullString `gorm:"-" json:"permission_names"`

	// Relations
	Permissions []Permission `gorm:"many2many:permissions_modules;" json:"permissions"`
}

// TableName specifies table name for GORM
func (PermissionGroup) TableName() string {
	return "permission_groups"
}
