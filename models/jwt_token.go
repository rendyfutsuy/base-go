package models

import (
	"time"

	"github.com/google/uuid"
)

// JWTToken represent the jwt token model
type JWTToken struct {
	UserId      uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"user_id" validate:"required"`
	AccessToken string     `gorm:"column:access_token;type:varchar(500);not null;primaryKey" json:"access_token" validate:"required"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies table name for GORM
func (JWTToken) TableName() string {
	return "jwt_tokens"
}
