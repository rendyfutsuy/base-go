package cobsubcob

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/dto"
)

type Usecase interface {
	InsertCategoryCobSubcob(cat []dto.CategoryJson, cob []dto.CobJson, createdById string) (err error)

	// cob
	GetCobByID(id string) (cobRes *models.Cob, err error)
	GetIndexCob(req *request.PageRequest) (cobs []models.Cob, total int, err error)
	GetAllCob() (cobs []models.Cob, err error)

	// sub cob
	GetSubcobByID(id string) (subcobRes *models.Subcob, err error)
	GetIndexSubcob(req *request.PageRequest) (subcobs []models.Subcob, total int, err error)
	GetAllSubcob() (subcobs []models.Subcob, err error)
}
