package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
)

// RequireActivatedUser blocks access if current authenticated user's IsFirstTimeLogin is true
func RequireActivatedUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userCtx := c.Get("user")
		currentUser, ok := userCtx.(models.User)
		if !ok {
			return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, "Unauthorized"))
		}
		if currentUser.IsFirstTimeLogin {
			return c.JSON(http.StatusForbidden, response.SetErrorResponse(http.StatusForbidden, constants.FirstTimeLoginErrorMessage))
		}
		return next(c)
	}
}
