package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type ReqCheckDuplicatedPermissionGroup struct {
	Name                      string    `json:"name" validate:"required"`
	ExcludedPermissionGroupId uuid.UUID `json:"excluded_role_info_id"`
}

type RespPermissionGroup struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value bool      `json:"value"`
}

type RespPermissionGroupIndex struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

type RespPermissionGroupDetail struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Permissions []string       `json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   utils.NullTime `json:"updated_at"`
}

// to get role info for compact use
func ToRespPermissionGroup(roleDb models.PermissionGroup) RespPermissionGroup {

	return RespPermissionGroup{
		ID:   roleDb.ID,
		Name: roleDb.Name,
	}

}

// for get role for index use
func ToRespPermissionGroupIndex(roleDb models.PermissionGroup) RespPermissionGroupIndex {

	return RespPermissionGroupIndex{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
	}

}

// to get role info with references
func ToRespPermissionGroupDetail(roleDb models.PermissionGroup) RespPermissionGroupDetail {

	permissions := make([]string, 0)
	for _, permission := range roleDb.PermissionNames {
		// Check if the permission is null or empty, continue to the next entry if it is
		if permission.Valid && permission.String != "" {
			// If the permission is valid and not empty, append it to the permissions slice
			permissions = append(permissions, permission.String)
		}
	}

	return RespPermissionGroupDetail{
		ID:          roleDb.ID,
		Name:        roleDb.Name,
		Permissions: permissions,
		CreatedAt:   roleDb.CreatedAt,
		UpdatedAt:   roleDb.UpdatedAt,
	}
}
