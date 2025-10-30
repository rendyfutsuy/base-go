package role_management

import (
	"context"

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
	CreateRole(ctx context.Context, roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (role *models.Role, err error)
	GetAllRole(ctx context.Context) (roles []models.Role, err error)
	GetIndexRole(ctx context.Context, req request.PageRequest) (roles []models.Role, total int, err error)
	UpdateRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error)
	SoftDeleteRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error)
	RoleNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedRole(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error)
	RoleNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedRoleOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error)

	CountRole(ctx context.Context) (count *int, err error)
	// ------------------------------------------------- role scope - END -----------------------------------------------------------

	// ------------------------------------------------- role assignment scope - BEGIN -----------------------------------------------------------
	ReAssignPermissionGroup(ctx context.Context, id uuid.UUID, permissionGroupReq dto.ToDBUpdatePermissionGroupAssignmentToRole) error
	GetTotalUser(ctx context.Context, id uuid.UUID) (total int, err error)
	GetPermissionFromRoleId(ctx context.Context, id uuid.UUID) (permissions []models.Permission, err error)
	GetPermissionGroupFromRoleId(ctx context.Context, id uuid.UUID) (permissionGroups []models.PermissionGroup, err error)
	AssignUsers(ctx context.Context, roleId uuid.UUID, userReq []uuid.UUID) error
	ReAssignPermissionsToPermissionGroup(ctx context.Context, id uuid.UUID, permissions []uuid.UUID) error
	GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error)

	// ------------------------------------------------- role assignment scope - END -----------------------------------------------------------

	// ------------------------------------------------- permission scope - BEGIN -----------------------------------------------------------
	// crud
	GetPermissionByID(ctx context.Context, id uuid.UUID) (permission *models.Permission, err error)
	GetAllPermission(ctx context.Context) (permissions []models.Permission, err error)
	GetIndexPermission(ctx context.Context, req request.PageRequest) (permissions []models.Permission, total int, err error)
	PermissionNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedPermission(ctx context.Context, name string, excludedId uuid.UUID) (permission *models.Permission, err error)

	CountPermission(ctx context.Context) (count *int, err error)
	// ------------------------------------------------- permission scope - END -----------------------------------------------------------

	// ------------------------------------------------- permission group scope - BEGIN -----------------------------------------------------------
	// crud
	GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (permissionGroup *models.PermissionGroup, err error)
	GetAllPermissionGroup(ctx context.Context) (permissionGroups []models.PermissionGroup, err error)
	GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (permissionGroups []models.PermissionGroup, total int, err error)
	PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedPermissionGroup(ctx context.Context, name string, excludedId uuid.UUID) (permissionGroup *models.PermissionGroup, err error)

	CountPermissionGroup(ctx context.Context) (count *int, err error)
	// ------------------------------------------------- permission group scope - END -----------------------------------------------------------
}
