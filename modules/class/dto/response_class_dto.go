package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespClass struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespClass(classDb models.Class) RespClass {

	return RespClass{
		ID:        classDb.ID,
		Name:      classDb.Name,
		Code:      classDb.Code,
		CreatedAt: classDb.CreatedAt,
		UpdatedAt: classDb.UpdatedAt,
	}

}
