package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type BackingSearchHelper struct{ request.SearchPredefineBase }

func (BackingSearchHelper) GetSearchColumns() []string          { return []string{"b.name", "b.backing_code"} }
func (BackingSearchHelper) GetSearchExistsSubqueries() []string { return []string{} }

var _ request.NeedSearchPredefine = BackingSearchHelper{}

func NewBackingSearchHelper() BackingSearchHelper {
	t := 0.75
	return BackingSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
