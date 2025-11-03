package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type ReqCheckDuplicatedPermission struct {
	Name                 string    `json:"name" validate:"required"`
	ExcludedPermissionId uuid.UUID `json:"excluded_role_info_id"`
}

type RespPermission struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RespPermissionIndex struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	TotalUser *int           `json:"total_user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

type RespPermissionDetail struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	TotalUser *int           `json:"total_user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

// to get role info for compact use
func ToRespPermission(roleDb models.Permission) RespPermission {

	return RespPermission{
		ID:   roleDb.ID,
		Name: roleDb.Name,
	}

}

// for get role for index use
func ToRespPermissionIndex(roleDb models.Permission) RespPermissionIndex {

	return RespPermissionIndex{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
	}

}

// to get role info with references
func ToRespPermissionDetail(roleDb models.Permission) RespPermissionDetail {

	return RespPermissionDetail{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
	}
}
