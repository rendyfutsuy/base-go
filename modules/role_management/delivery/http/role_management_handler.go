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
		"role.create", // store Role API
	}
	r.POST("/role", handler.CreateRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(storeRoles))

	// role show eligible permissions
	showRoles := []string{
		"role.view",   // index Role API
		"role.create", // store Role API
		"role.update", // update Role API
		"role.delete", // delete Role API
		"user.view",   // index User API
		"user.create", // store User API
		"user.update", // update User API
		"user.get",    // show User API
	}
	r.GET("/role", handler.GetIndexRole, middleware.RequireActivatedUser, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(showRoles))

	// role show eligible permissions
	allRoles := append(
		showRoles,  // same as show permissions assignment
		"role.all", // all Role API
	)
	r.GET("/role/all", handler.GetAllRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(allRoles))

	// role show eligible permissions
	showRoles = append(
		showRoles,  // same as show permissions assignment
		"role.get", // index Role API
	)
	r.GET("/role/:id", handler.GetRoleByID, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showRoles))

	// role update eligible permissions
	updateRoles := []string{
		"role.update", // update Role API
	}
	r.PUT("/role/:id", handler.UpdateRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(updateRoles))

	// role delete eligible permissions
	deleteRoles := []string{
		"role.delete", // delete Role API
	}
	r.DELETE("/role/:id", handler.DeleteRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(deleteRoles))

	// role show eligible permissions
	showRoles = append(
		showRoles,  // same as show permissions assignment
		"role.get", // index Role API
		"role.all", // all Role API
	)
	r.POST("/role/check-name", handler.GetDuplicatedRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showRoles))

	r.GET("/role/module-access", handler.GetAllPermissionGroupByModule, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showRoles))

	// role assignment
	// role re-assign permission groups eligible permissions
	reAssignPermissionGroups := []string{"role.re-assign-permission-groups"}
	r.PATCH("/role/:id/re-assign-permission-groups", handler.ReAssignPermissionByGroup, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(reAssignPermissionGroups))

	// role assign users eligible permissions
	assignUsers := []string{"role.assign-users"}
	r.PATCH("/role/:id/assign-users", handler.AssignUsersToRole, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(assignUsers))

	// 2025/11/04: unused - commented first
	// permission group scope
	// permission group index eligible permissions
	// indexPermissionGroups := []string{"permission-group.view", "role.update", "role.create"}
	// r.GET("/permission-group", handler.GetIndexPermissionGroup, middleware.RequireActivatedUser, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexPermissionGroups))

	// 2025/11/04: unused - commented first
	// permission group all eligible permissions
	// allPermissionGroups := append(
	// 	indexPermissionGroups,              // same as show permissions assignment
	// 	"permission-group.all",             // all permission group API
	// )
	// r.GET("/permission-group/all", handler.GetAllPermissionGroup, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(allPermissionGroups))

	// 2025/11/04: unused - commented first
	// permission group show eligible permissions
	// showPermissionGroups := append(
	// 	indexPermissionGroups,              // same as show permissions assignment
	// 	"permission-group.get",            // show permission group API
	// )
	// r.GET("/permission-group/:id", handler.GetPermissionGroupByID, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showPermissionGroups))

	// 2025/11/04: unused - commented first
	// permission group check name eligible permissions
	// checkNamePermissionGroups := append(
	// 	indexPermissionGroups,              // same as show permissions assignment
	// 	"permission-group.check-name",     // check name permission group API
	// 	"permission-group.all",            // all permission group API
	// 	"permission-group.get",           // show permission group API
	// )
	// r.GET("/permission-group/check-name", handler.GetDuplicatedPermissionGroup, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(checkNamePermissionGroups))

	// permission group all by module eligible permissions
	// allByModulePermissionGroups := append(
	// 	indexPermissionGroups,                    // same as show permissions assignment
	// 	"permission-group.all-by-module",        // all by module permission group API
	// 	"permission-group.all",                  // all permission group API
	// 	"permission-group.get",                  // show permission group API
	// )
	// r.GET("/permission-group/all/by-module", handler.GetAllPermissionGroupByModule, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(allByModulePermissionGroups))

	// 2025/11/04: unused - commented first
	// // get Current User Permissions
	// r.GET("/permission-group/all/my", handler.GetMyPermissions, middleware.RequireActivatedUser)

	// permission scope

	// 2025/11/04: unused - commented first
	// permission index eligible permissions
	// indexPermissions := []string{
	// 	"permission-group.all",   // all permission group API
	// 	"permission-group.view",  // index permission group API
	// 	"permission-group.get",   // show permission group API
	// 	"permission.view",        // index permission API
	// }
	// r.GET("/permission", handler.GetIndexPermission, middleware.RequireActivatedUser, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexPermissions))

	// 2025/11/04: unused - commented first
	// permission all eligible permissions
	// allPermissions := append(
	// 	indexPermissions,         // same as show permissions assignment
	// 	"permission.all",        // all permission API
	// )
	// r.GET("/permission/all", handler.GetAllPermission, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(allPermissions))

	// // permission show eligible permissions
	// showPermissions := append(
	// 	indexPermissions,         // same as show permissions assignment
	// 	"permission.get",        // show permission API
	// )
	// r.GET("/permission/:id", handler.GetPermissionByID, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showPermissions))

	// // permission check name eligible permissions
	// checkNamePermissions := append(
	// 	indexPermissions,         // same as show permissions assignment
	// 	"permission.check-name", // check name permission API
	// 	"permission.all",        // all permission API
	// 	"permission.get",        // show permission API
	// )
	// r.GET("/permission/check-name", handler.GetDuplicatedPermission, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(checkNamePermissions))
}
