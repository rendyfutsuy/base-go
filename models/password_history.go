package models

import (
	"time"

	"github.com/google/uuid"
)

// PasswordHistory represent the password history model
type PasswordHistory struct {
	UserId         uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"user_id" validate:"required"`
	HashedPassword string     `gorm:"column:hashed_password;type:varchar(255);not null;primaryKey" json:"hashed_password" validate:"required"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies table name for GORM
func (PasswordHistory) TableName() string {
	return "password_histories"
}
