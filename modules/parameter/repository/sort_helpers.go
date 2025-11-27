package repository

import "strings"

func normalizeParameterSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapParameterIndexSortColumn(sortBy string) string {
	normalized := normalizeParameterSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":           "p.id",
		"parameter_id": "p.id",
		"code":         "p.code",
		"name":         "p.name",
		"value":        "p.value",
		"type":         "p.type",
		"created_at":   "p.created_at",
		"updated_at":   "p.updated_at",
	}

	return mapping[normalized]
}
