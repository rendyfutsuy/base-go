package role

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	CreateRole(c echo.Context, req *dto.ReqCreateRole, authId uuid.UUID) (roleRes *models.Role, err error)
	GetRoleByID(id string) (role *models.Role, err error)
	GetAllRole() (roles []models.Role, err error)
	GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error)
	UpdateRole(id string, req *dto.ReqUpdateRole, authId uuid.UUID) (roleRes *models.Role, err error)
	SoftDeleteRole(id string, authId uuid.UUID) (roleRes *models.Role, err error)
}
