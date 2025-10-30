package models

import (
	"time"

	"github.com/google/uuid"
)

// ResetPasswordToken represent the reset password token model
type ResetPasswordToken struct {
	UserId      uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"user_id" validate:"required"`
	AccessToken string     `gorm:"column:access_token;type:varchar(500);not null;primaryKey" json:"access_token" validate:"required"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies table name for GORM
func (ResetPasswordToken) TableName() string {
	return "reset_password_tokens"
}
