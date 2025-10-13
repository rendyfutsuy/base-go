package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type RespRole struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RespRoleIndex struct {
	ID         uuid.UUID          `json:"id"`
	Name       string             `json:"name"`
	TotalUser  int                `json:"total_user"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  utils.NullTime     `json:"updated_at"`
	Modules    []utils.NullString `json:"modules"`
	Categories []utils.NullString `json:"units"`
}

type RespPermissionGroupRoleDetail struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Module string    `json:"module"`
}

type RespRoleDetail struct {
	ID               uuid.UUID                       `json:"id"`
	Name             string                          `json:"name"`
	TotalUser        int                             `json:"total_user"`
	PermissionGroups []RespPermissionGroupRoleDetail `json:"permission_groups"`
	Modules          []utils.NullString              `json:"modules"`
	CreatedAt        time.Time                       `json:"created_at"`
	UpdatedAt        utils.NullTime                  `json:"updated_at"`
	Description      utils.NullString                `json:"description"`
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
	Modules := make([]utils.NullString, 0)
	for _, Module := range roleDb.Modules {
		// append Module to array

		if Module.String != "" {
			Modules = append(Modules, Module)
		}
	}

	// mapping Cob to show at Role Detail
	Categories := make([]utils.NullString, 0)
	for _, Category := range roleDb.CategoryNames {
		// If the Category is valid and not empty, append it to the Categories slice
		if Category.String != "" {
			Categories = append(Categories, Category)
		}
	}

	return RespRoleIndex{
		ID:         roleDb.ID,
		Name:       roleDb.Name,
		TotalUser:  roleDb.TotalUser,
		Categories: Categories,
		CreatedAt:  roleDb.CreatedAt,
		UpdatedAt:  roleDb.UpdatedAt,
		Modules:    Modules,
	}

}

// to get role info with references
func ToRespRoleDetail(roleDb models.Role) RespRoleDetail {
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

	// mapping permission to show at Role Detail
	Modules := make([]utils.NullString, 0)
	for _, Module := range roleDb.Modules {
		// append Module to array

		if Module.String != "" {
			Modules = append(Modules, Module)
		}
	}

	return RespRoleDetail{
		ID:               roleDb.ID,
		Name:             roleDb.Name,
		TotalUser:        roleDb.TotalUser,
		PermissionGroups: PermissionGroups,
		Modules:          Modules,
		CreatedAt:        roleDb.CreatedAt,
		UpdatedAt:        roleDb.UpdatedAt,
		Description:      roleDb.Description,
	}
}
