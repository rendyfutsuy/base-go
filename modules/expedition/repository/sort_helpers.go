package repository

import "strings"

func normalizeExpeditionSortKey(sortBy string) string {
	sortBy = strings.TrimSpace(sortBy)
	if sortBy == "" {
		return ""
	}
	sortBy = strings.ReplaceAll(sortBy, "-", "_")
	sortBy = strings.ReplaceAll(sortBy, " ", "_")
	return strings.ToLower(sortBy)
}

func mapExpeditionIndexSortColumn(sortBy string) string {
	normalized := normalizeExpeditionSortKey(sortBy)
	if normalized == "" {
		return ""
	}

	mapping := map[string]string{
		"id":              "e.id",
		"expedition_id":   "e.id",
		"expedition_code": "e.expedition_code",
		"expedition_name": "e.expedition_name",
		"address":         "e.address",
		"phone_number":    "primary_phone_number",
		"telp_number":     "primary_telp_number",
		"created_at":      "e.created_at",
		"updated_at":      "e.updated_at",
	}

	return mapping[normalized]
}
