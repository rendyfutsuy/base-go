package request

import (
	"strings"

	"gorm.io/gorm"
)

type PageRequest struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
	Search    string `json:"search"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

func NewPageRequest(page, perPage int, search, sortBy, sortOrder string) *PageRequest {
	return &PageRequest{
		Page:      page,
		PerPage:   perPage,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// ApplySearchCondition applies search condition to a GORM query using ILIKE with OR conditions.
// If searchQuery is empty, the query is returned unchanged.
//
// Parameters:
//   - query: The GORM query builder
//   - searchQuery: The search string to match against
//   - searchColumns: List of column names to search in (e.g., []string{"usr.full_name", "usr.email"})
//
// Returns:
//   - *gorm.DB: The query with search condition applied (if searchQuery is not empty)
//
// Example:
//
//	query = ApplySearchCondition(query, "john", []string{"usr.full_name", "usr.email", "rl.name"})
//	// Results in: WHERE (usr.full_name ILIKE '%john%' OR usr.email ILIKE '%john%' OR rl.name ILIKE '%john%')
func ApplySearchCondition(query *gorm.DB, searchQuery string, searchColumns []string) *gorm.DB {
	if searchQuery == "" || len(searchColumns) == 0 {
		return query
	}

	// Build OR conditions for each column
	conditions := make([]string, len(searchColumns))
	args := make([]interface{}, len(searchColumns))
	searchPattern := "%" + searchQuery + "%"

	for i, column := range searchColumns {
		conditions[i] = column + " ILIKE ?"
		args[i] = searchPattern
	}

	// Combine with OR
	whereClause := "(" + strings.Join(conditions, " OR ") + ")"

	return query.Where(whereClause, args...)
}
