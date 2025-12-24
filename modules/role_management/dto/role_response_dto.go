package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type RespRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"role_name"`
}

type RespRoleIndex struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"role_name"`
	TotalUser int            `json:"total_user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
	Modules   []string       `json:"modules"`
	Deletable bool           `json:"deletable"`
}

type RespPermissionGroupRoleDetail struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Module string    `json:"module"`
}

type RespRoleDetail struct {
	ID          uuid.UUID                     `json:"id"`
	Name        string                        `json:"role_name"`
	TotalUser   int                           `json:"total_user"`
	Modules     []RespPermissionGroupByModule `json:"modules"`
	CreatedAt   time.Time                     `json:"created_at"`
	UpdatedAt   utils.NullTime                `json:"updated_at"`
	Description string                        `json:"description"`
	Deletable   bool                          `json:"deletable"`
}

// to get role info for compact use
func ToRespRole(roleDb models.Role) RespRole {

	return RespRole{
		ID:   roleDb.ID,
		Name: roleDb.Name,
	}

}

// for get role for index use
func ToRespRoleIndex(roleDb models.Role) RespRoleIndex {

	// mapping permission to show at Role Detail
	Modules := make([]string, 0)
	for _, Module := range roleDb.Modules {
		// append Module to array

		if Module.String != "" {
			Modules = append(Modules, Module.String)
		}
	}

	// mapping Cob to show at Role Detail
	Categories := make([]string, 0)
	for _, Category := range roleDb.CategoryNames {
		// If the Category is valid and not empty, append it to the Categories slice
		if Category.String != "" {
			Categories = append(Categories, Category.String)
		}
	}

	// mapping permission group to show at Role Detail
	deletable := roleDb.Deletable

	// if total user > 0, set deletable to false
	if roleDb.TotalUser > 0 {
		deletable = false
	}

	return RespRoleIndex{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		TotalUser: roleDb.TotalUser,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
		Modules:   Modules,
		Deletable: deletable,
	}

}

// to get role info with references
func ToRespRoleDetail(roleDb models.Role, modules []RespPermissionGroupByModule) RespRoleDetail {
	// mapping permission group to show at Role Detail
	PermissionGroups := make([]RespPermissionGroupRoleDetail, 0)
	for _, PermissionGroup := range roleDb.PermissionGroups {
		// If the PermissionGroup is valid and not empty, append it to the PermissionGroups slice
		PermissionGroups = append(PermissionGroups, RespPermissionGroupRoleDetail{
			ID:     PermissionGroup.ID,
			Name:   PermissionGroup.Name,
			Module: PermissionGroup.Module.String,
		})
	}

	// mapping permission group to show at Role Detail
	deletable := roleDb.Deletable

	// if total user > 0, set deletable to false
	if roleDb.TotalUser > 0 {
		deletable = false
	}

	return RespRoleDetail{
		ID:          roleDb.ID,
		Name:        roleDb.Name,
		TotalUser:   roleDb.TotalUser,
		Modules:     modules,
		CreatedAt:   roleDb.CreatedAt,
		UpdatedAt:   roleDb.UpdatedAt,
		Description: roleDb.Description.String,
		Deletable:   deletable,
	}
}
