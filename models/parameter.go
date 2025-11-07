package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Parameter represents parameter table
type Parameter struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	Code        string         `gorm:"column:code;type:varchar(255);unique;not null" json:"code"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	Value       *string        `gorm:"column:value;type:varchar(255)" json:"value"`
	Type        *string        `gorm:"column:type;type:varchar(255)" json:"type"`
	Description *string        `gorm:"column:description;type:text" json:"desc"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Parameter) TableName() string {
	return "parameter"
}
