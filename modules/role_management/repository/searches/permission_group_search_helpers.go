package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type PermissionGroupSearchHelper struct{ request.SearchPredefineBase }

func (PermissionGroupSearchHelper) GetSearchColumns() []string {
	return []string{"permission_group.name"}
}
func (PermissionGroupSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = PermissionGroupSearchHelper{}

func NewPermissionGroupSearchHelper() PermissionGroupSearchHelper {
	t := 0.50
	return PermissionGroupSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
