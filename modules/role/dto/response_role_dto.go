package dto

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/utils"
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
