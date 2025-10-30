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
	GetRoleByID(id string) (role *models.Role, err error)
	GetAllRole() (role_infos []models.Role, err error)
	GetIndexRole(req request.PageRequest) (role_infos []models.Role, total int, err error)
	UpdateRole(id string, req *dto.ReqUpdateRole, authId string) (roleRes *models.Role, err error)
	SoftDeleteRole(id string, authId string) (roleRes *models.Role, err error)
	RoleNameIsNotDuplicated(name string, id uuid.UUID) (roleRes *models.Role, err error)
	MyPermissionsByUserToken(c echo.Context, token string) (role *models.Role, err error)

	// role assignment scope
	ReAssignPermissionByGroup(roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (roleRes *models.Role, err error)
	AssignUsersToRole(roleId string, req *dto.ReqUpdateAssignUsersToRole) (roleRes *models.Role, err error)

	// permission group scope
	GetPermissionGroupByID(id string) (role *models.PermissionGroup, err error)
	GetAllPermissionGroup() (role_infos []models.PermissionGroup, err error)
	GetIndexPermissionGroup(req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error)
	PermissionGroupNameIsNotDuplicated(name string, id uuid.UUID) (roleRes *models.PermissionGroup, err error)

	// permission scope
	GetPermissionByID(id string) (role *models.Permission, err error)
	GetAllPermission() (role_infos []models.Permission, err error)
	GetIndexPermission(req request.PageRequest) (role_infos []models.Permission, total int, err error)
	PermissionNameIsNotDuplicated(name string, id uuid.UUID) (roleRes *models.Permission, err error)
}
