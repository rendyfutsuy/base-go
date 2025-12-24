package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type ParameterSearchHelper struct{ request.SearchPredefineBase }

func (ParameterSearchHelper) GetSearchColumns() []string {
	return []string{"p.name", "p.code"}
}
func (ParameterSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = ParameterSearchHelper{}

func NewParameterSearchHelper() ParameterSearchHelper {
	t := 0.75
	return ParameterSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
