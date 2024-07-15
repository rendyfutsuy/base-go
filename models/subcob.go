package models

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type Subcob struct {
	ID                      uuid.UUID        `json:"id"`
	CategoryID              uuid.UUID        `json:"category_id"`
	CobID                   uuid.UUID        `json:"cob_id"`
	Name                    string           `json:"name"`
	Code                    string           `json:"code"`
	Forms                   utils.NullString `json:"forms"`
	ActiveDate              utils.NullTime   `json:"active_date"`
	IsHiddenFromFacultative utils.NullBool   `json:"is_hidden_from_facultative"`
	IsInactive              utils.NullBool   `json:"is_inactive"`
	IsFromWebCrediit        utils.NullBool   `json:"is_from_web_credit"`
	TimestampBase
}
