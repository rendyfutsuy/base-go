package models

import (
	"time"

	"github.com/google/uuid"
)

// JWTToken represent the jwt token model
type JWTToken struct {
	UserId           uuid.UUID  `gorm:"type:uuid;not null;primaryKey" json:"user_id" validate:"required"`
	AccessToken      string     `gorm:"column:access_token;type:varchar(500);not null;primaryKey" json:"access_token" validate:"required"`
	AccessJTI        string     `gorm:"column:access_jti;type:varchar(255)" json:"access_jti"`
	RefreshToken     string     `gorm:"column:refresh_token;type:varchar(500);not null" json:"refresh_token"`
	RefreshJTI       string     `gorm:"column:refresh_jti;type:varchar(255)" json:"refresh_jti"`
	RefreshExpiresAt time.Time  `gorm:"column:refresh_expires_at" json:"refresh_expires_at"`
	IsUsed           bool       `gorm:"column:is_used;default:false" json:"is_used"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName specifies table name for GORM
func (JWTToken) TableName() string {
	return "jwt_tokens"
}
