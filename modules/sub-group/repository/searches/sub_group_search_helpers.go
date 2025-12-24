package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type SubGroupSearchHelper struct{ request.SearchPredefineBase }

func (SubGroupSearchHelper) GetSearchColumns() []string {
	return []string{"sg.name", "sg.subgroup_code", "gg.name"}
}
func (SubGroupSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = SubGroupSearchHelper{}

func NewSubGroupSearchHelper() SubGroupSearchHelper {
	return SubGroupSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: nil}}
}
