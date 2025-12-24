package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
// implement search on city scope -- BEGIN
type CitySearchHelper struct{ request.SearchPredefineBase }

func (CitySearchHelper) GetSearchColumns() []string {
	return []string{"c.name", "c.area_code"}
}
func (CitySearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = CitySearchHelper{}

func NewCitySearchHelper() CitySearchHelper {
	t := 0.40
	return CitySearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}

// implement search on city scope -- END
