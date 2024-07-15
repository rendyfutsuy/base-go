package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespContractor struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Address   string         `json:"address"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespContractor(contractorDb models.Contractor) RespContractor {
	return RespContractor{
		ID:        contractorDb.ID,
		Name:      contractorDb.Name,
		Code:      contractorDb.Code,
		Address:   contractorDb.Address,
		CreatedAt: contractorDb.CreatedAt,
		UpdatedAt: contractorDb.UpdatedAt,
	}
}
