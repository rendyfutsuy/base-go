package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rendyfutsuy/base-go/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/role_management"
)

type IMiddlewarePermission interface {
	PermissionValidation(args []string) echo.MiddlewareFunc
}

type MiddlewarePermission struct {
	authRepository           auth.Repository
	roleManagementRepository role_management.Repository
}

func NewMiddlewarePermission(authRepository auth.Repository, roleManagementRepository role_management.Repository) IMiddlewarePermission {
	return &MiddlewarePermission{
		authRepository:           authRepository,
		roleManagementRepository: roleManagementRepository,
	}
}

func (a *MiddlewarePermission) PermissionValidation(args []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			authorization := c.Request().Header.Get("Authorization")

			if authorization == "" {
				return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Unauthorized: No Authorization header"})
			}

			tokenParts := strings.Split(authorization, "Bearer ")
			if len(tokenParts) != 2 {
				return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Unauthorized: Invalid Token format"})
			}

			tokenString := tokenParts[1]
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Unauthorized: Token not provided"})
			}

			// get user data from token
			user, err := a.getUserData(ctx, tokenString)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Unauthorized: Token invalid"})
			}

			// fetch permissions based on role user's has
			// get user data from token
			permissions, err := a.getUserPermissions(ctx, user.RoleId)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, GeneralResponse{Message: "Unauthorized: Unable to fetch permissions"})
			}

			// compare if there match between permissions and requiredPermissions
			if !a.assertUserHaveRequiredPermissions(permissions, args) {
				return c.JSON(http.StatusForbidden, GeneralResponse{Message: "Forbidden: Insufficient permissions"})
			}

			return next(c)
		}
	}
}

func (a *MiddlewarePermission) getUserData(ctx context.Context, token string) (models.User, error) {
	// get user data from current token
	return a.authRepository.GetUserByAccessToken(ctx, token)
}

func (a *MiddlewarePermission) getUserPermissions(ctx context.Context, roleUid uuid.UUID) ([]string, error) {
	// get permissions from user's role
	role, err := a.roleManagementRepository.GetRoleByID(ctx, roleUid)
	if err != nil {
		return nil, err
	}

	// mapped permission to to array string
	var permissions []string
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Name)
	}
	return permissions, nil
}

func (a *MiddlewarePermission) assertUserHaveRequiredPermissions(userPermissions []string, requiredPermissions []string) bool {
	permissionSet := make(map[string]bool)
	// assign user permissions to compared permissions
	for _, p := range userPermissions {
		permissionSet[p] = true
	}

	// assert permission required existed in permissions owns by user
	for _, rp := range requiredPermissions {
		// assert at least 1 permission from required permissions exists in user permissions
		if permissionSet[rp] {
			return true
		}
	}
	return false
}
