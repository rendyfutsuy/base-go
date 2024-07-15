package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespRole struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespRole(roleDb models.Role) RespRole {

	return RespRole{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
	}

}
