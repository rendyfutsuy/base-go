package models

import (
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Token     string     `gorm:"type:varchar(255);not null" json:"token"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OTP) TableName() string {
	return "otps"
}
