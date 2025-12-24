package role_management

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

type Usecase interface {
	// role scope
	CreateRole(ctx context.Context, req *dto.ReqCreateRole, authId string) (roleRes *models.Role, err error)
	GetRoleByID(ctx context.Context, id string) (role *models.Role, err error)
	GetAllRole(ctx context.Context) (role_infos []models.Role, err error)
	GetIndexRole(ctx context.Context, req request.PageRequest) (role_infos []models.Role, total int, err error)
	UpdateRole(ctx context.Context, id string, req *dto.ReqUpdateRole, authId string) (roleRes *models.Role, err error)
	SoftDeleteRole(ctx context.Context, id string, authId string) (roleRes *models.Role, err error)
	RoleNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (roleRes *models.Role, err error)
	MyPermissionsByUserToken(ctx context.Context, token string) (role *models.Role, err error)

	// role assignment scope
	ReAssignPermissionByGroup(ctx context.Context, roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (roleRes *models.Role, err error)
	AssignUsersToRole(ctx context.Context, roleId string, req *dto.ReqUpdateAssignUsersToRole) (roleRes *models.Role, err error)

	// permission group scope
	GetPermissionGroupByID(ctx context.Context, id string) (role *models.PermissionGroup, err error)
	GetAllPermissionGroup(ctx context.Context) (role_infos []models.PermissionGroup, err error)
	GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error)
	PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (roleRes *models.PermissionGroup, err error)

	// permission scope
	GetPermissionByID(ctx context.Context, id string) (role *models.Permission, err error)
	GetAllPermission(ctx context.Context) (role_infos []models.Permission, err error)
	GetIndexPermission(ctx context.Context, req request.PageRequest) (role_infos []models.Permission, total int, err error)
	PermissionNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (roleRes *models.Permission, err error)
}
