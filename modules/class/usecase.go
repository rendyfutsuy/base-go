package class

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/dto"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	CreateClass(c echo.Context, req *dto.ReqCreateClass, authId string) (classRes *models.Class, err error)
	GetClassByID(id string) (class *models.Class, err error)
	GetAllClass() (classes []models.Class, err error)
	GetIndexClass(req request.PageRequest) (classes []models.Class, total int, err error)
	UpdateClass(id string, req *dto.ReqUpdateClass, authId string) (classRes *models.Class, err error)
	SoftDeleteClass(id string, authId string) (classRes *models.Class, err error)
}
