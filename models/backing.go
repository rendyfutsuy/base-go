package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Backing represents backings table
type Backing struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	TypeID     uuid.UUID      `gorm:"column:type_id;type:uuid;not null" json:"type_id" validate:"required"`
	BackingCode string        `gorm:"column:backing_code;type:varchar(255);unique;not null" json:"backing_code"`
	Name       string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	CreatedBy  string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy  string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy  *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`

	// Relations
	Type *Type `gorm:"foreignKey:TypeID" json:"type,omitempty"`

	// Read-only fields from join (not stored in database)
	TypeName      string    `gorm:"column:type_name;<-:false" json:"type_name"`
	SubgroupID    uuid.UUID `gorm:"column:subgroup_id;<-:false" json:"subgroup_id"`
	SubgroupName  string    `gorm:"column:subgroup_name;<-:false" json:"subgroup_name"`
	GoodsGroupID  uuid.UUID `gorm:"column:goods_group_id;<-:false" json:"goods_group_id"`
	GroupName     string    `gorm:"column:group_name;<-:false" json:"group_name"`
}

func (Backing) TableName() string {
	return "backings"
}

