package repository

import "strings"

func normalizeGroupSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapGroupIndexSortColumn(sortBy string) string {
	normalized := normalizeGroupSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":         "gg.id",
		"group_id":   "gg.id",
		"group_code": "gg.group_code",
		"name":       "gg.name", // Natural sorting enabled for this column
		"created_at": "gg.created_at",
		"updated_at": "gg.updated_at",
	}

	return mapping[normalized]
}
