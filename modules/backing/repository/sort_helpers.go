package repository

import "strings"

func normalizeBackingSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapBackingIndexSortColumn(sortBy string) string {
	normalized := normalizeBackingSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":            "b.id",
		"backing_id":    "b.id",
		"type_id":       "b.type_id",
		"backing_code":  "b.backing_code",
		"name":          "b.name",
		"type_name":     "type_name",
		"subgroup_name": "subgroup_name",
		"group_name":    "group_name",
		"created_at":    "b.created_at",
		"updated_at":    "b.updated_at",
	}

	return mapping[normalized]
}
