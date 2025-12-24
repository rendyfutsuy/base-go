package dto

import (
	"github.com/rendyfutsuy/base-go/models"
)

type RespPermissionGroupByModule struct {
	Name             string                `json:"name"`
	PermissionGroups []RespPermissionGroup `json:"accesses"`
}

// to get role info for compact use
func ToRespPermissionGroupByModule(roleDb models.PermissionGroup) RespPermissionGroupByModule {

	return RespPermissionGroupByModule{
		Name: roleDb.Module.String,
	}

}
