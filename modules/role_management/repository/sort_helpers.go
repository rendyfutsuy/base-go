package repository

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
func mapRoleIndexSortColumn(sortBy string) string {
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

// mapPermissionIndexSortColumn maps DTO index fields to sortable columns for the permissions index.
// The returned value should match the AllowedColumns used with ApplyPagination (without prefix).
func mapPermissionIndexSortColumn(sortBy string) string {
	normalized := normalizeSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":                    "id",
		"permission.id":         "id",
		"name":                  "name",
		"permission.name":       "name",
		"created_at":            "created_at",
		"permission.created_at": "created_at",
		"updated_at":            "updated_at",
		"permission.updated_at": "updated_at",
		"deleted_at":            "deleted_at",
		"permission.deleted_at": "deleted_at",
	}

	return mapping[normalized]
}

// mapPermissionGroupIndexSortColumn maps DTO index fields to sortable columns for the permission groups index.
// The returned value should match the AllowedColumns used with ApplyPagination (without prefix).
func mapPermissionGroupIndexSortColumn(sortBy string) string {
	normalized := normalizeSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":                          "id",
		"permission_group.id":         "id",
		"name":                        "name",
		"permission_group.name":       "name",
		"module":                      "module",
		"permission_group.module":     "module",
		"created_at":                  "created_at",
		"permission_group.created_at": "created_at",
		"updated_at":                  "updated_at",
		"permission_group.updated_at": "updated_at",
		"deleted_at":                  "deleted_at",
		"permission_group.deleted_at": "deleted_at",
	}

	return mapping[normalized]
}
