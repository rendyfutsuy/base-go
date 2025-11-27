package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/modules/user_management"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserManagementHandler struct {
	UserUseCase          user_management.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewUserManagementHandler(e *echo.Echo, us user_management.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	handler := &UserManagementHandler{
		UserUseCase:          us,
		validator:            validator.New(),
		mwPageRequest:        mwP,
		middlewareAuth:       auth,
		middlewarePermission: middlewarePermission,
	}

	r := e.Group("v1/user-management")

	r.Use(handler.middlewareAuth.AuthorizationCheck)

	// Permissions
	permissionToCreate := []string{"user.create"}
	permissionToUpdate := []string{"user.update"}
	permissionToDelete := []string{"user.delete"}

	// user scope
	r.POST("/user", handler.CreateUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(permissionToCreate))

	// user index eligible permission
	indexUser := []string{
		"user.view",     // access user index
		"user.update",   // access user update
		"user.block",    // access user block
		"user.activate", // access user activate
		"user.create",   // access user store
	}
	r.GET("/user", handler.GetIndexUser, middleware.RequireActivatedUser, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexUser))

	// user show eligible permission
	allUser := append(indexUser, "user.all")
	r.GET("/user/all", handler.GetAllUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(allUser))

	// user show eligible permission
	showUser := append(indexUser, "user.get")
	r.GET("/user/:id", handler.GetUserByID, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(showUser))

	// user update
	r.PUT("/user/:id", handler.UpdateUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(permissionToUpdate))

	// user delete
	r.DELETE("/user/:id", handler.DeleteUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(permissionToDelete))

	// user check name	eligible permission
	checkUserName := append(showUser, "user.check-name")
	r.POST("/user/check-name", handler.GetDuplicatedUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(checkUserName))

	// 2025/11/04: unused - commented first
	// user block
	// r.PATCH("/user/:id/block", handler.BlockUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation([]string{"user.block"}))

	// 2025/11/11: unused - commented first
	// user activate
	// r.PATCH("/user/:id/assign-status", handler.ActivateUser, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation([]string{"user.activate"}))

	// user update password
	// Allow password update without RequireActivatedUser so user can activate themselves
	r.PATCH("/user/:id/password", handler.UpdateUserPassword, handler.middlewarePermission.PermissionValidation([]string{"user.update-password"}))

	// user password confirmation
	// Allow password confirmation without RequireActivatedUser
	r.POST("/user/password-confirmation", handler.ConfirmCurrentUserPassword)

	// user import from Excel
	r.GET("/user/import/template", handler.DownloadUserImportTemplate, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(permissionToCreate))
	r.POST("/user/import", handler.ImportUsersFromExcel, middleware.RequireActivatedUser, handler.middlewarePermission.PermissionValidation(permissionToCreate))

}
