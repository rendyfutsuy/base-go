package dto

import (
	"github.com/google/uuid"
)

type ReqCheckDuplicatedRole struct {
	Name           string    `json:"role_name" validate:"required"`
	ExcludedRoleId uuid.UUID `json:"excluded_role_info_id"`
}

type ReqCreateRole struct {
	Name             string      `form:"role_name" json:"role_name" validate:"required,max=80"`
	Description      string      `form:"description" json:"description"`
	PermissionGroups []uuid.UUID `form:"accesses" json:"accesses" validate:"required,min=1"`
}

func (r *ReqCreateRole) ToDBCreateRole(code, authId string) ToDBCreateRole {
	return ToDBCreateRole{
		Name:             r.Name,
		Description:      r.Description,
		PermissionGroups: r.PermissionGroups,
	}
}

type ToDBCreateRole struct {
	Name             string      `json:"role_name"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"accesses"`
}
