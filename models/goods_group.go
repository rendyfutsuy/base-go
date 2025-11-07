package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// GoodsGroup represents goods_group table
type GoodsGroup struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" validate:"required"`
	GroupCode string         `gorm:"column:group_code;type:varchar(255);unique;not null" json:"group_code"`
	Name      string         `gorm:"column:name;type:varchar(255);not null;uniqueIndex" json:"name" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (GoodsGroup) TableName() string {
    return "goods_group"
}


