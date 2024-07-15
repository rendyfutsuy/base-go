package carriage

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	// crud
	CreateCarriage(c echo.Context, req *dto.ReqCreateCarriage, authId string) (carriageRes *models.Carriage, err error)
	GetCarriageByID(id string) (carriage *models.Carriage, err error)
	GetAllCarriage() (carriagees []models.Carriage, err error)
	GetIndexCarriage(req request.PageRequest) (carriagees []models.Carriage, total int, err error)
	UpdateCarriage(id string, req *dto.ReqUpdateCarriage, authId string) (carriageRes *models.Carriage, err error)
	SoftDeleteCarriage(id string, authId string) (carriageRes *models.Carriage, err error)
}
