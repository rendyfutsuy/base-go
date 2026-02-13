package dto

import "github.com/google/uuid"

type ReqUpdateRole struct {
	Name             string      `form:"role_name" json:"role_name" validate:"required,max=80"`
	Description      string      `form:"description" json:"description"`
	PermissionGroups []uuid.UUID `form:"accesses" json:"accesses" validate:"required,min=1"`
}

func (r *ReqUpdateRole) ToDBUpdateRole(authId string) ToDBUpdateRole {
	return ToDBUpdateRole{
		Name:        r.Name,
		Description: r.Description,
	}
}

type ToDBUpdateRole struct {
	Name             string      `json:"role_name"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"accesses"`
}
