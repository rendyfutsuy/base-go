package role

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	models "github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/modules/role/dto"
)

type Usecase interface {
	CreateRole(c echo.Context, req *dto.ReqCreateRole, authId uuid.UUID) (roleRes *models.Role, err error)
	GetRoleByID(id string) (role *models.Role, err error)
	GetAllRole() (roles []models.Role, err error)
	GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error)
	UpdateRole(id string, req *dto.ReqUpdateRole, authId uuid.UUID) (roleRes *models.Role, err error)
	SoftDeleteRole(id string, authId uuid.UUID) (roleRes *models.Role, err error)
}
