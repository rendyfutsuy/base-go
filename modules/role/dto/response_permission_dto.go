package dto

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/utils"
)

type RespPermission struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespPermission(roleDb models.Role) RespPermission {

	return RespPermission{
		ID:        roleDb.ID,
		Name:      roleDb.Name,
		CreatedAt: roleDb.CreatedAt,
		UpdatedAt: roleDb.UpdatedAt,
	}

}
