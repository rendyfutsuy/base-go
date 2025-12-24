package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
// for regency make new struct to handle province, city, district, subdistrict separately
// this is to avoid code duplication and conflict
// for threshold setting on all regency is set as nil.
// so threshold on each regency will be use the default threshold

// implement search on province scope -- BEGIN
type ProvinceSearchHelper struct{ request.SearchPredefineBase }

func (ProvinceSearchHelper) GetSearchColumns() []string {
	return []string{"p.name"}
}
func (ProvinceSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = ProvinceSearchHelper{}

func NewProvinceSearchHelper() ProvinceSearchHelper {
	t := 0.33
	return ProvinceSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}

// implement search on province scope -- END
