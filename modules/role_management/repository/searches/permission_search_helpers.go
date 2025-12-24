package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type PermissionSearchHelper struct{ request.SearchPredefineBase }

func (PermissionSearchHelper) GetSearchColumns() []string {
	return []string{"permission.name"}
}
func (PermissionSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = PermissionSearchHelper{}

func NewPermissionSearchHelper() PermissionSearchHelper {
	t := 0.55
	return PermissionSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
