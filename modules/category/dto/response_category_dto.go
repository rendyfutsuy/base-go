package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespCategory struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	CreatedAt   utils.NullTime `json:"created_at"`
	UpdatedAt   utils.NullTime `json:"updated_at"`
}

func ToRespCategory(categoryDb models.Category) RespCategory {

	return RespCategory{
		ID:          categoryDb.ID,
		Name:        categoryDb.Name,
		Code:        categoryDb.Code,
		Description: categoryDb.Description,
		CreatedAt:   categoryDb.CreatedAt,
		UpdatedAt:   categoryDb.UpdatedAt,
	}

}
