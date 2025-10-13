package role_management

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// ------------------------------------------------- role scope - BEGIN -----------------------------------------------------------
	// crud
	CreateRole(roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error)
	GetRoleByID(id uuid.UUID) (role *models.Role, err error)
	GetAllRole() (roles []models.Role, err error)
	GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error)
	UpdateRole(id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error)
	SoftDeleteRole(id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error)
	RoleNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedRole(name string, excludedId uuid.UUID) (role *models.Role, err error)
	RoleNameIsNotDuplicatedOnSoftDeleted(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedRoleOnSoftDeleted(name string, excludedId uuid.UUID) (role *models.Role, err error)

	CountRole() (count *int, err error)
	// ------------------------------------------------- role scope - END -----------------------------------------------------------

	// ------------------------------------------------- role assignment scope - BEGIN -----------------------------------------------------------
	ReAssignPermissionGroup(id uuid.UUID, permissionGroupReq dto.ToDBUpdatePermissionGroupAssignmentToRole) error
	GetTotalUser(id uuid.UUID) (total int, err error)
	GetPermissionFromRoleId(id uuid.UUID) (permissions []models.Permission, err error)
	GetPermissionGroupFromRoleId(id uuid.UUID) (permissionGroups []models.PermissionGroup, err error)
	AssignUsers(roleId uuid.UUID, userReq []uuid.UUID) error
	ReAssignPermissionsToPermissionGroup(id uuid.UUID, permissions []uuid.UUID) error
	GetUserByID(id uuid.UUID) (user *models.User, err error)

	// ------------------------------------------------- role assignment scope - END -----------------------------------------------------------

	// ------------------------------------------------- permission scope - BEGIN -----------------------------------------------------------
	// crud
	GetPermissionByID(id uuid.UUID) (permission *models.Permission, err error)
	GetAllPermission() (permissions []models.Permission, err error)
	GetIndexPermission(req request.PageRequest) (permissions []models.Permission, total int, err error)
	PermissionNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedPermission(name string, excludedId uuid.UUID) (permission *models.Permission, err error)

	CountPermission() (count *int, err error)
	// ------------------------------------------------- permission scope - END -----------------------------------------------------------

	// ------------------------------------------------- permission group scope - BEGIN -----------------------------------------------------------
	// crud
	GetPermissionGroupByID(id uuid.UUID) (permissionGroup *models.PermissionGroup, err error)
	GetAllPermissionGroup() (permissionGroups []models.PermissionGroup, err error)
	GetIndexPermissionGroup(req request.PageRequest) (permissionGroups []models.PermissionGroup, total int, err error)
	PermissionGroupNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedPermissionGroup(name string, excludedId uuid.UUID) (permissionGroup *models.PermissionGroup, err error)

	CountPermissionGroup() (count *int, err error)
	// ------------------------------------------------- permission group scope - END -----------------------------------------------------------
}
