package role

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helper/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role/dto"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// crud
	CreateRole(roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error)
	GetRoleByID(id uuid.UUID) (role *models.Role, err error)
	GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error)
	GetAllRole() (roles []models.Role, err error)
	UpdateRole(id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error)
	SoftDeleteRole(id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error)

	// general
	CountRole() (count *int, err error)
}
