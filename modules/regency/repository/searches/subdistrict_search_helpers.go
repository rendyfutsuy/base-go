package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
// implement search on subdistrict scope -- BEGIN
type SubdistrictSearchHelper struct{ request.SearchPredefineBase }

func (SubdistrictSearchHelper) GetSearchColumns() []string {
	return []string{"s.name"}
}
func (SubdistrictSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = SubdistrictSearchHelper{}

func NewSubdistrictSearchHelper() SubdistrictSearchHelper {
	t := 0.50
	return SubdistrictSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}

// implement search on subdistrict scope -- END
