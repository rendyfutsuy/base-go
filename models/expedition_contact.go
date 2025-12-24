package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExpeditionContact represents expedition_contacts table
type ExpeditionContact struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	ExpeditionID uuid.UUID      `gorm:"column:expedition_id;type:uuid;not null" json:"expedition_id" validate:"required"`
	PhoneType   string         `gorm:"column:phone_type;type:varchar(50);not null" json:"phone_type" validate:"required"` // telp / hp
	PhoneNumber string         `gorm:"column:phone_number;type:varchar(50);not null" json:"phone_number" validate:"required"`
	AreaCode    *string        `gorm:"column:area_code;type:varchar(255)" json:"area_code"`
	IsPrimary   bool           `gorm:"column:is_primary;default:false" json:"is_primary"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy   string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	UpdatedBy   string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy   *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`
}

func (ExpeditionContact) TableName() string {
	return "expedition_contacts"
}

