package role_management

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

type Usecase interface {
	// role scope
	CreateRole(c echo.Context, req *dto.ReqCreateRole, authId string) (roleRes *models.Role, err error)
	GetRoleByID(c echo.Context, id string) (role *models.Role, err error)
	GetAllRole(c echo.Context) (role_infos []models.Role, err error)
	GetIndexRole(c echo.Context, req request.PageRequest) (role_infos []models.Role, total int, err error)
	UpdateRole(c echo.Context, id string, req *dto.ReqUpdateRole, authId string) (roleRes *models.Role, err error)
	SoftDeleteRole(c echo.Context, id string, authId string) (roleRes *models.Role, err error)
	RoleNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (roleRes *models.Role, err error)
	MyPermissionsByUserToken(c echo.Context, token string) (role *models.Role, err error)

	// role assignment scope
	ReAssignPermissionByGroup(c echo.Context, roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (roleRes *models.Role, err error)
	AssignUsersToRole(c echo.Context, roleId string, req *dto.ReqUpdateAssignUsersToRole) (roleRes *models.Role, err error)

	// permission group scope
	GetPermissionGroupByID(c echo.Context, id string) (role *models.PermissionGroup, err error)
	GetAllPermissionGroup(c echo.Context) (role_infos []models.PermissionGroup, err error)
	GetIndexPermissionGroup(c echo.Context, req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error)
	PermissionGroupNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (roleRes *models.PermissionGroup, err error)

	// permission scope
	GetPermissionByID(c echo.Context, id string) (role *models.Permission, err error)
	GetAllPermission(c echo.Context) (role_infos []models.Permission, err error)
	GetIndexPermission(c echo.Context, req request.PageRequest) (role_infos []models.Permission, total int, err error)
	PermissionNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (roleRes *models.Permission, err error)
}
