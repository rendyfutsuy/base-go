package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespCarriage struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespCarriage(carriageDb models.Carriage) RespCarriage {
	return RespCarriage{
		ID:        carriageDb.ID,
		Name:      carriageDb.Name,
		Code:      carriageDb.Code,
		CreatedAt: carriageDb.CreatedAt,
		UpdatedAt: carriageDb.UpdatedAt,
	}
}
