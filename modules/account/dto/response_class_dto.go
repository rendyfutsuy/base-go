package dto

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type RespAccount struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	CreatedAt utils.NullTime `json:"created_at"`
	UpdatedAt utils.NullTime `json:"updated_at"`
}

func ToRespAccount(accountDb models.Account) RespAccount {

	return RespAccount{
		ID:        accountDb.ID,
		Name:      accountDb.Name,
		Code:      accountDb.Code,
		CreatedAt: accountDb.CreatedAt,
		UpdatedAt: accountDb.UpdatedAt,
	}

}
