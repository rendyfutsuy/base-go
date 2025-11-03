package dto

import (
	"github.com/google/uuid"
)

type ReqCheckDuplicatedRole struct {
	Name           string    `json:"name" validate:"required"`
	ExcludedRoleId uuid.UUID `json:"excluded_role_info_id"`
}

type ReqCreateRole struct {
	Name             string      `json:"name" validate:"required,max=80"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"permission_groups" validate:"required,min=1"`
}

func (r *ReqCreateRole) ToDBCreateRole(code, authId string) ToDBCreateRole {
	return ToDBCreateRole{
		Name:             r.Name,
		Description:      r.Description,
		PermissionGroups: r.PermissionGroups,
	}
}

type ToDBCreateRole struct {
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"permission_groups"`
	Cobs             []uuid.UUID `json:"cobs"`
	Categories       []uuid.UUID `json:"categories"`
}
