package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Expedition represents expeditions table
type Expedition struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	ExpeditionCode string         `gorm:"column:expedition_code;type:varchar(255);unique;not null" json:"expedition_code"`
	ExpeditionName string         `gorm:"column:expedition_name;type:varchar(255)" json:"expedition_name"`
	Address        string         `gorm:"column:address;type:varchar(255)" json:"address"`
	Notes          *string        `gorm:"column:notes;type:text" json:"notes"`
	CreatedAt      time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy      string         `gorm:"column:created_by;type:varchar(255)" json:"created_by"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	UpdatedBy      string         `gorm:"column:updated_by;type:varchar(255)" json:"updated_by"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	DeletedBy      *string        `gorm:"column:deleted_by;type:varchar(255)" json:"deleted_by"`

	// Fetched mutators for primary contacts (from joins)
	PrimaryTelpNumber  *string `gorm:"column:primary_telp_number;<-:false" json:"primary_telp_number"`
	PrimaryPhoneNumber *string `gorm:"column:primary_phone_number;<-:false" json:"primary_phone_number"`
}

func (Expedition) TableName() string {
	return "expeditions"
}

