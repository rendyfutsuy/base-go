package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/modules/role_management"
)

type ResponseError struct {
	Message string `json:"message"`
}

type RoleManagementHandler struct {
	RoleUseCase          role_management.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewRoleManagementHandler(e *echo.Echo, us role_management.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	handler := &RoleManagementHandler{
		RoleUseCase:          us,
		validator:            validator.New(),
		mwPageRequest:        mwP,
		middlewareAuth:       auth,
		middlewarePermission: middlewarePermission,
	}

	r := e.Group("v1/role-management")

	r.Use(handler.middlewareAuth.AuthorizationCheck)

	// role scope assignment
	// role store eligible permissions
	storeRoles := []string{
		"api.role-management.role.store", // store Role API
	}
	r.POST("/role", handler.CreateRole, handler.middlewarePermission.PermissionValidation(storeRoles))

	// role show eligible permissions
	showRoles := []string{
		"api.role-management.role.index",  // index Role API
		"api.role-management.role.store",  // store Role API
		"api.role-management.role.update", // update Role API
		"api.role-management.role.delete", // delete Role API
		"api.user-management.user.index",  // index User API
		"api.user-management.user.store",  // store User API
		"api.user-management.user.update", // update User API
		"api.user-management.user.show",   // show User API
	}
	r.GET("/role", handler.GetIndexRole, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(showRoles))

	// role show eligible permissions
	allRoles := append(
		showRoles,                      // same as show permissions assignment
		"api.role-management.role.all", // all Role API
	)
	r.GET("/role/all", handler.GetAllRole, handler.middlewarePermission.PermissionValidation(allRoles))

	// role show eligible permissions
	showRoles = append(
		showRoles,                       // same as show permissions assignment
		"api.role-management.role.show", // index Role API
	)
	r.GET("/role/:id", handler.GetRoleByID, handler.middlewarePermission.PermissionValidation(showRoles))

	// role update eligible permissions
	updateRoles := []string{
		"api.role-management.role.update", // update Role API
	}
	r.PUT("/role/:id", handler.UpdateRole, handler.middlewarePermission.PermissionValidation(updateRoles))

	// role delete eligible permissions
	deleteRoles := []string{
		"api.role-management.role.delete", // delete Role API
	}
	r.DELETE("/role/:id", handler.DeleteRole, handler.middlewarePermission.PermissionValidation(deleteRoles))

	// role show eligible permissions
	showRoles = append(
		showRoles,                       // same as show permissions assignment
		"api.role-management.role.show", // index Role API
		"api.role-management.role.all",  // all Role API
	)
	r.GET("/role/check-name", handler.GetDuplicatedRole, handler.middlewarePermission.PermissionValidation(showRoles))

	// role assignment
	// role re-assign permission groups eligible permissions
	reAssignPermissionGroups := []string{"api.role-management.role.re-assign-permission-groups"}
	r.PATCH("/role/:id/re-assign-permission-groups", handler.ReAssignPermissionByGroup, handler.middlewarePermission.PermissionValidation(reAssignPermissionGroups))

	// role assign users eligible permissions
	assignUsers := []string{"api.role-management.role.assign-users"}
	r.PATCH("/role/:id/assign-users", handler.AssignUsersToRole, handler.middlewarePermission.PermissionValidation(assignUsers))

	// permission group scope
	// permission group index eligible permissions
	indexPermissionGroups := []string{"api.role-management.permission-groups.index", "api.role-management.role.update", "api.role-management.role.store"}
	r.GET("/permission-group", handler.GetIndexPermissionGroup, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexPermissionGroups))

	// permission group all eligible permissions
	allPermissionGroups := append(
		indexPermissionGroups,                       // same as show permissions assignment
		"api.role-management.permission-groups.all", // all permission group API
	)
	r.GET("/permission-group/all", handler.GetAllPermissionGroup, handler.middlewarePermission.PermissionValidation(allPermissionGroups))

	// permission group show eligible permissions
	showPermissionGroups := append(
		indexPermissionGroups,                        // same as show permissions assignment
		"api.role-management.permission-groups.show", // show permission group API
	)
	r.GET("/permission-group/:id", handler.GetPermissionGroupByID, handler.middlewarePermission.PermissionValidation(showPermissionGroups))

	// permission group check name eligible permissions
	checkNamePermissionGroups := append(
		indexPermissionGroups,                              // same as show permissions assignment
		"api.role-management.permission-groups.check-name", // check name permission group API
		"api.role-management.permission-groups.all",        // all permission group API
		"api.role-management.permission-groups.show",       // show permission group API
	)
	r.GET("/permission-group/check-name", handler.GetDuplicatedPermissionGroup, handler.middlewarePermission.PermissionValidation(checkNamePermissionGroups))

	// permission group all by module eligible permissions
	allByModulePermissionGroups := append(
		indexPermissionGroups, // same as show permissions assignment
		"api.role-management.permission-groups.all-by-module", // all by module permission group API
		"api.role-management.permission-groups.all",           // all permission group API
		"api.role-management.permission-groups.show",          // show permission group API
	)
	r.GET("/permission-group/all/by-module", handler.GetAllPermissionGroupByModule, handler.middlewarePermission.PermissionValidation(allByModulePermissionGroups))

	// get Current User Permissions
	r.GET("/permission-group/all/my", handler.GetMyPermissions)

	// permission scope

	// permission index eligible permissions
	indexPermissions := []string{
		"api.role-management.permission-groups.all",   // all permission group API
		"api.role-management.permission-groups.index", // index permission group API
		"api.role-management.permission-groups.show",  // show permission group API
		"api.role-management.permission.index",        // index permission API
	}
	r.GET("/permission", handler.GetIndexPermission, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexPermissions))

	// permission all eligible permissions
	allPermissions := append(
		indexPermissions,                     // same as show permissions assignment
		"api.role-management.permission.all", // all permission API
	)
	r.GET("/permission/all", handler.GetAllPermission, handler.middlewarePermission.PermissionValidation(allPermissions))

	// permission show eligible permissions
	showPermissions := append(
		indexPermissions,                      // same as show permissions assignment
		"api.role-management.permission.show", // show permission API
	)
	r.GET("/permission/:id", handler.GetPermissionByID, handler.middlewarePermission.PermissionValidation(showPermissions))

	// permission check name eligible permissions
	checkNamePermissions := append(
		indexPermissions, // same as show permissions assignment
		"api.role-management.permission.check-name", // check name permission API
		"api.role-management.permission.all",        // all permission API
		"api.role-management.permission.show",       // show permission API
	)
	r.GET("/permission/check-name", handler.GetDuplicatedPermission, handler.middlewarePermission.PermissionValidation(checkNamePermissions))
}
