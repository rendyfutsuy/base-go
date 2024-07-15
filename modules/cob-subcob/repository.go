package cobsubcob

import (
	"database/sql"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// transaction
	StartTransaction() (*sql.Tx, error)

	// cob
	CreateCob(trx *sql.Tx, cobReq dto.ToDBCreateCob) (cobRes *models.Cob, err error)
	GetCobByID(id uuid.UUID) (cobRes *models.Cob, err error)
	GetIndexCob(req request.PageRequest) (cobs []models.Cob, total int, err error)
	GetAllCob() (cobs []models.Cob, err error)

	// subcob
	CreateSubcob(trx *sql.Tx, subcobReq dto.ToDBCreateSubcob) (subcobRes *models.Subcob, err error)
	GetSubcobByID(id uuid.UUID) (subcobRes *models.Subcob, err error)
	GetIndexSubcob(req request.PageRequest) (Subcobs []models.Subcob, total int, err error)
	GetAllSubcob() (subcobs []models.Subcob, err error)
}
