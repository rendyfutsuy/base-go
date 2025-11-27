package usecase

import (
	"strings"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
)

func isSuperAdminRoleName(name string) bool {
	return strings.EqualFold(name, constants.AuthRoleSuperAdmin)
}

func isRestrictedUserPermissionGroup(permissionGroup *models.PermissionGroup) bool {
	if permissionGroup == nil {
		return false
	}

	if !permissionGroup.Module.Valid {
		return false
	}

	if !strings.EqualFold(permissionGroup.Module.String, constants.UserPermissionModuleName) {
		return false
	}

	return strings.EqualFold(permissionGroup.Name, constants.UserPermissionNameCreate) ||
		strings.EqualFold(permissionGroup.Name, constants.UserPermissionNameDelete)
}
