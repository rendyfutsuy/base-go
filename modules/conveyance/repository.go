package conveyance

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// migration
	// CreateTable(sqlFilePath string) (err error)

	// crud
	CreateConveyance(conveyanceReq dto.ToDBCreateConveyance) (conveyanceRes *models.Conveyance, err error)
	GetConveyanceByID(id uuid.UUID) (conveyance *models.Conveyance, err error)
	GetIndexConveyance(req request.PageRequest) (conveyancees []models.Conveyance, total int, err error)
	GetAllConveyance() (conveyancees []models.Conveyance, err error)
	UpdateConveyance(id uuid.UUID, conveyanceReq dto.ToDBUpdateConveyance) (conveyanceRes *models.Conveyance, err error)
	SoftDeleteConveyance(id uuid.UUID, conveyanceReq dto.ToDBDeleteConveyance) (conveyanceRes *models.Conveyance, err error)

	// general
	CountConveyance() (count *int, err error)
}
