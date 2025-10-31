package dto

import "github.com/google/uuid"

type ReqUpdateRole struct {
	Name             string      `json:"name" validate:"required,max=80"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"permission_groups" validate:"required,min=1"`
	Cobs             []uuid.UUID `json:"cobs" validate:"required,min=1"`
}

func (r *ReqUpdateRole) ToDBUpdateRole(authId string) ToDBUpdateRole {
	return ToDBUpdateRole{
		Name:        r.Name,
		Description: r.Description,
		Cobs:        r.Cobs,
	}
}

type ToDBUpdateRole struct {
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	PermissionGroups []uuid.UUID `json:"permission_groups"`
	Cobs             []uuid.UUID `json:"cobs"`
}
