package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type ReqConfirmationUserPassword struct {
	Password string `json:"password" validate:"required"`
}

type RespUser struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"name"`
}

type RespUserIndex struct {
	ID           uuid.UUID        `json:"id"`
	FullName     string           `json:"name"`
	Email        string           `json:"email"`
	IsBlocked    bool             `json:"is_blocked"`
	RoleName     string           `json:"role_name"`
	IsActive     bool             `json:"is_active"`
	ActiveStatus utils.NullString `json:"active_status"`
	Gender       string           `json:"gender"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type ReqUserIndexFilter struct {
	RoleIds  []uuid.UUID `query:"role_ids"`
	RoleName string      `query:"role_name"`
}

type RespPermissionGroupUserDetail struct {
	ID           uuid.UUID        `json:"id"`
	FullName     string           `json:"name"`
	Email        string           `json:"email"`
	IsBlocked    bool             `json:"is_blocked"`
	RoleName     string           `json:"role_name"`
	RoleId       uuid.UUID        `json:"role_id"`
	IsActive     bool             `json:"is_active"`
	ActiveStatus utils.NullString `json:"active_status"`
}

type RespUserDetail struct {
	ID           uuid.UUID        `json:"id"`
	FullName     string           `json:"name"`
	Email        string           `json:"email"`
	IsBlocked    bool             `json:"is_blocked"`
	RoleId       uuid.UUID        `json:"role_id"`
	RoleName     string           `json:"role_name"`
	IsActive     bool             `json:"is_active"`
	ActiveStatus utils.NullString `json:"active_status"`
	Gender       string           `json:"gender"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// to get role info for compact use
func ToRespUser(userDb models.User) RespUser {

	return RespUser{
		ID:       userDb.ID,
		FullName: userDb.FullName,
	}

}

// for get role for index use
func ToRespUserIndex(userDb models.User) RespUserIndex {

	return RespUserIndex{
		ID:           userDb.ID,
		FullName:     userDb.FullName,
		Email:        userDb.Email,
		IsBlocked:    userDb.IsBlocked,
		IsActive:     userDb.IsActive,
		ActiveStatus: userDb.ActiveStatus,
		RoleName:     userDb.RoleName,
		Gender:       userDb.Gender,
		CreatedAt:    userDb.CreatedAt,
		UpdatedAt:    userDb.UpdatedAt,
	}

}

// to get role info with references
func ToRespUserDetail(userDb models.User) RespUserDetail {

	return RespUserDetail{
		ID:           userDb.ID,
		FullName:     userDb.FullName,
		Email:        userDb.Email,
		IsBlocked:    userDb.IsBlocked,
		RoleId:       userDb.RoleId,
		RoleName:     userDb.RoleName,
		IsActive:     userDb.IsActive,
		ActiveStatus: userDb.ActiveStatus,
		Gender:       userDb.Gender,
		CreatedAt:    userDb.CreatedAt,
		UpdatedAt:    userDb.UpdatedAt,
	}
}
