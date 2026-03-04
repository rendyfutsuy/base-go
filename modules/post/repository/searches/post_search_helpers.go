package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

type PostSearchHelper struct{ request.SearchPredefineBase }

func (PostSearchHelper) GetSearchColumns() []string {
	return []string{
		"c.title",
		"c.short_description",
	}
}

func (PostSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{}
}

var _ request.NeedSearchPredefine = PostSearchHelper{}

func NewPostSearchHelper() PostSearchHelper {
	t := 0.75
	return PostSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
