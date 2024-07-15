package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespConveyance struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Type      string         `json:"type"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespConveyance(conveyanceDb models.Conveyance) RespConveyance {
	return RespConveyance{
		ID:        conveyanceDb.ID,
		Name:      conveyanceDb.Name,
		Code:      conveyanceDb.Code,
		Type:      conveyanceDb.Type,
		CreatedAt: conveyanceDb.CreatedAt,
		UpdatedAt: conveyanceDb.UpdatedAt,
	}
}
