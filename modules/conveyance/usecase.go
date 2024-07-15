package conveyance

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	// crud
	CreateConveyance(c echo.Context, req *dto.ReqCreateConveyance, authId string) (conveyanceRes *models.Conveyance, err error)
	GetConveyanceByID(id string) (conveyance *models.Conveyance, err error)
	GetAllConveyance() (conveyancees []models.Conveyance, err error)
	GetIndexConveyance(req request.PageRequest) (conveyancees []models.Conveyance, total int, err error)
	UpdateConveyance(id string, req *dto.ReqUpdateConveyance, authId string) (conveyanceRes *models.Conveyance, err error)
	SoftDeleteConveyance(id string, authId string) (conveyanceRes *models.Conveyance, err error)
}
