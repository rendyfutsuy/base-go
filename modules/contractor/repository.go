package contractor

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// migration
	// CreateTable(sqlFilePath string) (err error)

	// crud
	CreateContractor(contractorReq dto.ToDBCreateContractor) (contractorRes *models.Contractor, err error)
	GetContractorByID(id uuid.UUID) (contractor *models.Contractor, err error)
	GetIndexContractor(req request.PageRequest) (contractores []models.Contractor, total int, err error)
	GetAllContractor() (contractores []models.Contractor, err error)
	UpdateContractor(id uuid.UUID, contractorReq dto.ToDBUpdateContractor) (contractorRes *models.Contractor, err error)
	SoftDeleteContractor(id uuid.UUID, contractorReq dto.ToDBDeleteContractor) (contractorRes *models.Contractor, err error)

	// general
	CountContractor() (count *int, err error)
}
