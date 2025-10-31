package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

// User represent the user model
type User struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id" validate:"required"`
	RoleId            uuid.UUID      `gorm:"type:uuid;not null" json:"role_id" validate:"required"`
	FullName          string         `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" validate:"required"`
	Username          string         `gorm:"column:username;type:varchar(100)" json:"username" validate:"required"`
	Email             string         `gorm:"column:email;type:varchar(255);not null;uniqueIndex" json:"email" validate:"required,email"`
	Password          string         `gorm:"column:password;type:varchar(255);not null" json:"password" validate:"required"`
	IsActive          bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	PasswordExpiredAt time.Time      `gorm:"column:password_expired_at" json:"password_expired_at"`
	Gender            string         `gorm:"column:gender;type:varchar(20)" json:"gender"`
	Counter           int            `gorm:"column:counter;default:0" json:"counter"`

	// mutator - not stored in DB
	ActiveStatus utils.NullString `gorm:"-" json:"active_status"`
	IsBlocked    bool             `gorm:"-" json:"is_blocked"`
	RoleName     string           `gorm:"-" json:"role_name"`
}

// TableName specifies table name for GORM
func (User) TableName() string {
	return "users"
}
