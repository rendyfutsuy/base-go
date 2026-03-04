package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null" json:"name"`
	FilePath    *string        `gorm:"column:file_path;type:text" json:"file_path"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (File) TableName() string {
	return "files"
}
