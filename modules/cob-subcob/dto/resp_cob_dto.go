package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespCob struct {
	ID         uuid.UUID      `json:"id"`
	CategoryID uuid.UUID      `json:"category_id"`
	Name       string         `json:"name"`
	Code       string         `json:"code"`
	CreatedAt  utils.NullTime `json:"created_at"`
	UpdatedAt  utils.NullTime `json:"updated_at"`
}

func ToRespCob(cobDb models.Cob) RespCob {

	return RespCob{
		ID:         cobDb.ID,
		CategoryID: cobDb.CategoryID,
		Name:       cobDb.Name,
		Code:       cobDb.Code,
		CreatedAt:  cobDb.CreatedAt,
		UpdatedAt:  cobDb.UpdatedAt,
	}

}
