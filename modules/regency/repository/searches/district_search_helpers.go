package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
// for regency make new struct to handle province, city, district, subdistrict separately
// this is to avoid code duplication and conflict
// for threshold setting on all regency is set as nil.
// so threshold on each regency will be use the default threshold
// implement search on district scope -- BEGIN
type DistrictSearchHelper struct{ request.SearchPredefineBase }

func (DistrictSearchHelper) GetSearchColumns() []string {
	return []string{"d.name"}
}
func (DistrictSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = DistrictSearchHelper{}

func NewDistrictSearchHelper() DistrictSearchHelper {
	t := 0.50
	return DistrictSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}

// implement search on district scope -- END
