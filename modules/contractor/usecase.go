package contractor

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	// crud
	CreateContractor(c echo.Context, req *dto.ReqCreateContractor, authId string) (contractorRes *models.Contractor, err error)
	GetContractorByID(id string) (contractor *models.Contractor, err error)
	GetAllContractor() (contractores []models.Contractor, err error)
	GetIndexContractor(req request.PageRequest) (contractores []models.Contractor, total int, err error)
	UpdateContractor(id string, req *dto.ReqUpdateContractor, authId string) (contractorRes *models.Contractor, err error)
	SoftDeleteContractor(id string, authId string) (contractorRes *models.Contractor, err error)
}
