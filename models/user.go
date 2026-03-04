package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/utils"
	utilsServices "github.com/rendyfutsuy/base-go/utils/services"
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
	Nik               string         `gorm:"column:nik;type:varchar(80)" json:"nik"`
	IsActive          bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	VerifiedAt        *time.Time     `gorm:"column:verified_at" json:"verified_at"`
	PasswordExpiredAt time.Time      `gorm:"column:password_expired_at" json:"password_expired_at"`
	Gender            string         `gorm:"column:gender;type:varchar(20)" json:"gender"`
	Counter           int            `gorm:"column:counter;default:0" json:"counter"`
	IsFirstTimeLogin  bool           `gorm:"column:is_first_time_login" json:"is_first_time_login"`
	Deletable         bool           `gorm:"column:deletable;default:true;not null" json:"deletable"`
	// Files relation (pivot)
	Files []File `gorm:"many2many:files_to_module;joinForeignKey:ID;joinReferences:FileID" json:"-"`

	// mutator - not stored in DB
	ActiveStatus     utils.NullString `gorm:"column:active_status;<-:false" json:"active_status"` // Read-only: used for fetch, ignored on insert/update
	IsBlocked        bool             `gorm:"column:is_blocked;<-:false" json:"is_blocked"`       // Read-only: used for fetch, ignored on insert/update
	RoleName         string           `gorm:"column:role_name;<-:false" json:"role_name"`         // Read-only: used for fetch, ignored on insert/update
	AvatarURL        *string          `gorm:"column:avatar_url;<-:false" json:"avatar_url"`       // Read-only from pivot
	Permissions      []string         `gorm:"-" json:"permissions"`
	PermissionGroups []string         `gorm:"-" json:"permission_groups"`
	Modules          []string         `gorm:"-" json:"modules"`
}

// TableName specifies table name for GORM
func (User) TableName() string {
	return "users"
}

func (user User) GetAvatarURL() string {
	if user.AvatarURL == nil {
		return ""
	}
	presignedURL, err := utilsServices.GeneratePresignedURL(*user.AvatarURL)
	if err != nil {
		return ""
	}

	return presignedURL
}

// AfterFind computes AvatarURL from files_to_module pivot with type = "avatar"
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	if u == nil || u.ID == uuid.Nil {
		return nil
	}
	var filePath *string
	err = tx.Table("files_to_module ftm").
		Select("f.file_path").
		Joins("JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL").
		Where("ftm.module_type = ? AND ftm.module_id = ? AND ftm.type = ?", constants.ModuleTypeUser, u.ID, constants.FileTypeAvatar).
		Order("ftm.created_at DESC").
		Limit(1).
		Scan(&filePath).Error
	if err != nil {
		return nil // do not block read on error
	}
	if filePath != nil {
		u.AvatarURL = filePath
	}
	return nil
}
