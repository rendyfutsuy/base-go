package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type GroupSearchHelper struct{ request.SearchPredefineBase }

func (GroupSearchHelper) GetSearchColumns() []string {
	return []string{"gg.name", "gg.group_code"}
}
func (GroupSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = GroupSearchHelper{}

func NewGroupSearchHelper() GroupSearchHelper {
	t := 0.50
	return GroupSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
