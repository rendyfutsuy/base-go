package sorts

import "strings"

// normalizeSortKey converts incoming sort_by values into a comparable snake_case form.
// It trims whitespace, replaces spaces and hyphens with underscores, and lowers the case.
func normalizeSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

// mapRoleIndexSortColumn maps DTO index fields to actual sortable columns for the roles index.
// It returns empty string if the provided key is not recognized.
func MapRoleIndexSortColumn(sortBy string) string {
	normalized := normalizeSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":              "role.id",
		"role.id":         "role.id",
		"role_name":       "role.name",
		"role.name":       "role.name",
		"total_user":      "total_user",
		"created_at":      "role.created_at",
		"role.created_at": "role.created_at",
		"updated_at":      "role.updated_at",
		"role.updated_at": "role.updated_at",
		"modules":         "modules",
		"role.modules":    "modules",
	}

	return mapping[normalized]
}
