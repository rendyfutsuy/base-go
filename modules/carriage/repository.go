package carriage

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// migration
	// CreateTable(sqlFilePath string) (err error)

	// crud
	CreateCarriage(carriageReq dto.ToDBCreateCarriage) (carriageRes *models.Carriage, err error)
	GetCarriageByID(id uuid.UUID) (carriage *models.Carriage, err error)
	GetIndexCarriage(req request.PageRequest) (carriagees []models.Carriage, total int, err error)
	GetAllCarriage() (carriagees []models.Carriage, err error)
	UpdateCarriage(id uuid.UUID, carriageReq dto.ToDBUpdateCarriage) (carriageRes *models.Carriage, err error)
	SoftDeleteCarriage(id uuid.UUID, carriageReq dto.ToDBDeleteCarriage) (carriageRes *models.Carriage, err error)

	// general
	CountCarriage() (count *int, err error)
}
