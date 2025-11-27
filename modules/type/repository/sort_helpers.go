package repository

import "strings"

func normalizeTypeSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapTypeIndexSortColumn(sortBy string) string {
	normalized := normalizeTypeSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":               "t.id",
		"type_id":          "t.id",
		"subgroup_id":      "t.subgroup_id",
		"type_code":        "t.type_code",
		"name":             "t.name",
		"subgroup_name":    "subgroup_name",
		"goods_group_name": "goods_group_name",
		"group_name":       "goods_group_name",
		"created_at":       "t.created_at",
		"updated_at":       "t.updated_at",
	}

	return mapping[normalized]
}
