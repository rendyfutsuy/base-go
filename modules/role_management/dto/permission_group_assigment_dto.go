package dto

import (
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
)

type RespPermissionGroupByModule struct {
	Name             utils.NullString      `json:"name"`
	PermissionGroups []RespPermissionGroup `json:"permission_groups"`
}

// to get role info for compact use
func ToRespPermissionGroupByModule(roleDb models.PermissionGroup) RespPermissionGroupByModule {

	return RespPermissionGroupByModule{
		Name: roleDb.Module,
	}

}
