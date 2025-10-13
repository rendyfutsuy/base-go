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

	// user scope
	r.POST("/user", handler.CreateUser, handler.middlewarePermission.PermissionValidation([]string{"api.user-management.user.store"}))

	// user index eligible permission
	indexUser := []string{
		"api.user-management.user.index",              // access user index
		"api.user.management.user.update",             // access user update
		"api.user-management.user.block",              // access user block
		"api.user-management.user.activate",           // access user activate
		"api.user-management.user.store",              // access user store
		"api.facultative.integration.user-management", // integration with Facultative
		"api.treaty.integration.user-management",      // integration with treaty
		"api.claim.integration.user-management",       // integration with claim
	}
	r.GET("/user", handler.GetIndexUser, handler.mwPageRequest.PageRequestCtx, handler.middlewarePermission.PermissionValidation(indexUser))

	// user show eligible permission
	allUser := append(indexUser, "api.user-management.user.all")
	r.GET("/user/all", handler.GetAllUser, handler.middlewarePermission.PermissionValidation(allUser))

	// user show eligible permission
	showUser := append(indexUser, "api.user-management.user.show")
	r.GET("/user/:id", handler.GetUserByID, handler.middlewarePermission.PermissionValidation(showUser))

	// user update
	r.PUT("/user/:id", handler.UpdateUser, handler.middlewarePermission.PermissionValidation([]string{"api.user-management.user.update"}))

	// user check name	eligible permission
	checkUserName := append(showUser, "api.user-management.user.check-name")
	r.GET("/user/check-name", handler.GetDuplicatedUser, handler.middlewarePermission.PermissionValidation(checkUserName))

	// user block
	r.PATCH("/user/:id/block", handler.BlockUser, handler.middlewarePermission.PermissionValidation([]string{"api.user-management.user.block"}))

	// user activate
	r.PATCH("/user/:id/assign-status", handler.ActivateUser, handler.middlewarePermission.PermissionValidation([]string{"api.user-management.user.activate"}))

	// user update password
	r.PATCH("/user/:id/password", handler.UpdateUserPassword, handler.middlewarePermission.PermissionValidation([]string{"api.user-management.user.update-password"}))

	// user password confirmation
	r.POST("/user/password-confirmation", handler.ConfirmCurrentUserPassword)

}
