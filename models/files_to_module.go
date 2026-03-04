package models

import (
	"time"

	"github.com/google/uuid"
)

type FilesToModule struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	FileID     uuid.UUID `gorm:"column:file_id;type:uuid;not null" json:"file_id"`
	ModuleType string    `gorm:"column:module_type;type:varchar(255);not null" json:"module_type"`
	ModuleID   uuid.UUID `gorm:"column:module_id;type:uuid;not null" json:"module_id"`
	Type       string    `gorm:"column:type;type:varchar(255)" json:"type"`
	CreatedAt  time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (FilesToModule) TableName() string {
	return "files_to_module"
}
