package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Type represents types table
type Type struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	SubgroupID  uuid.UUID      `gorm:"column:subgroup_id;type:uuid;not null" json:"subgroup_id" validate:"required"`
	TypeCode    string         `gorm:"column:type_code;type:varchar(255);not null;unique" json:"type_code"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy   string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	UpdatedBy   string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy   *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`

	// Read-only field from join (not stored in database)
	SubgroupName   string    `gorm:"column:subgroup_name;<-:false" json:"subgroup_name"`
	GoodsGroupID   uuid.UUID `gorm:"column:goods_group_id;<-:false" json:"goods_group_id"`
	GoodsGroupName string    `gorm:"column:goods_group_name;<-:false" json:"goods_group_name"`
}

func (Type) TableName() string {
	return "types"
}

