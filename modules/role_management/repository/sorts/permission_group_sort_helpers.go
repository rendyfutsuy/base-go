package sorts

// mapPermissionGroupIndexSortColumn maps DTO index fields to sortable columns for the permission groups index.
// The returned value should match the AllowedColumns used with ApplyPagination (without prefix).
func MapPermissionGroupIndexSortColumn(sortBy string) string {
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
