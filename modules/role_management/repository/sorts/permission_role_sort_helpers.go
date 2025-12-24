package sorts

// mapPermissionIndexSortColumn maps DTO index fields to sortable columns for the permissions index.
// The returned value should match the AllowedColumns used with ApplyPagination (without prefix).
func MapPermissionIndexSortColumn(sortBy string) string {
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
