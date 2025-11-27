package repository

import "strings"

func normalizeSubGroupSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapSubGroupIndexSortColumn(sortBy string) string {
	normalized := normalizeSubGroupSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":               "sg.id",
		"subgroup_id":      "sg.id",
		"goods_group_id":   "sg.goods_group_id",
		"subgroup_code":    "sg.subgroup_code",
		"name":             "sg.name",
		"goods_group_name": "goods_group_name",
		"created_at":       "sg.created_at",
		"updated_at":       "sg.updated_at",
	}

	return mapping[normalized]
}
