package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Group represents groups table
type Group struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" validate:"required"`
	GroupCode string         `gorm:"column:group_code;type:varchar(255);unique;not null" json:"group_code"`
	Name      string         `gorm:"column:name;type:varchar(255);not null;uniqueIndex" json:"name" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	UpdatedBy string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`

	// Read-only field from query (not stored in database)
	Deletable bool `gorm:"column:deletable;<-:false" json:"deletable"`
}

func (Group) TableName() string {
	return "groups"
}
