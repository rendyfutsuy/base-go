package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type RoleSearchHelper struct{ request.SearchPredefineBase }

func (RoleSearchHelper) GetSearchColumns() []string {
	return []string{"role.name"}
}
func (RoleSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = RoleSearchHelper{}

func NewRoleSearchHelper() RoleSearchHelper {
	t := 0.70
	return RoleSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
