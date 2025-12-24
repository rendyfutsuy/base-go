package searches

import "github.com/rendyfutsuy/base-go/helpers/request"

// initialize, value for search and map the function & variable need for it
type TypeSearchHelper struct{ request.SearchPredefineBase }

func (TypeSearchHelper) GetSearchColumns() []string {
	return []string{
		"t.type_code",
		"t.name",
	}
}

func (TypeSearchHelper) GetSearchExistsSubqueries() []string {
	return []string{
		"EXISTS (SELECT 1 FROM sub_groups sg WHERE sg.id = t.subgroup_id AND sg.deleted_at IS NULL AND REPLACE(sg.name, ' ', '') ILIKE ?)",
		"EXISTS (SELECT 1 FROM groups gg JOIN sub_groups sg2 ON gg.id = sg2.groups_id WHERE sg2.id = t.subgroup_id AND gg.deleted_at IS NULL AND sg2.deleted_at IS NULL AND REPLACE(gg.name, ' ', '') ILIKE ?)",
	}
}

var _ request.NeedSearchPredefine = TypeSearchHelper{}

func NewTypeSearchHelper() TypeSearchHelper {
	t := 0.75
	return TypeSearchHelper{SearchPredefineBase: request.SearchPredefineBase{Threshold: &t}}
}
