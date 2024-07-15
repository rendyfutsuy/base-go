package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespSubcob struct {
	ID         uuid.UUID      `json:"id"`
	CategoryID uuid.UUID      `json:"category_id"`
	CobID      uuid.UUID      `json:"cob_id"`
	Name       string         `json:"name"`
	Code       string         `json:"code"`
	CreatedAt  utils.NullTime `json:"created_at"`
	UpdatedAt  utils.NullTime `json:"updated_at"`
}

func ToRespSubcob(cobDb models.Subcob) RespSubcob {

	return RespSubcob{
		ID:         cobDb.ID,
		CategoryID: cobDb.CategoryID,
		CobID:      cobDb.CobID,
		Name:       cobDb.Name,
		Code:       cobDb.Code,
		CreatedAt:  cobDb.CreatedAt,
		UpdatedAt:  cobDb.UpdatedAt,
	}

}
