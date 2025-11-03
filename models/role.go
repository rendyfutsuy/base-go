package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
)

// Role represent the role model
type Role struct {
	ID          uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	Name        string           `gorm:"type:varchar(255);not null;uniqueIndex" json:"name" validate:"required"`
	Deletable   bool             `gorm:"column:deletable;default:true;not null" json:"deletable"`
	CreatedAt   time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   utils.NullTime   `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   utils.NullTime   `gorm:"column:deleted_at;index" json:"deleted_at"`
	Description utils.NullString `gorm:"column:description;type:text" json:"description"`

	// Computed/virtual fields - not stored in DB
	TotalUser            int                `gorm:"-" json:"total_user"`
	Modules              []utils.NullString `gorm:"-" json:"modules"`
	CategoryNames        []utils.NullString `gorm:"-" json:"category_names"`
	PermissionGroupNames []utils.NullString `gorm:"-" json:"permission_group_names"`
	PermissionGroupIds   []uuid.UUID        `gorm:"-" json:"permission_group_ids"`

	// Relations
	Users            []User            `gorm:"foreignKey:RoleId" json:"users"`
	Permissions      []Permission      `gorm:"many2many:permissions_modules;" json:"permissions"`
	PermissionGroups []PermissionGroup `gorm:"many2many:modules_roles;" json:"permission_groups"`
}

// TableName specifies table name for GORM
func (Role) TableName() string {
	return "roles"
}
