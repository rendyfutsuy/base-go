package request

import (
	"fmt"
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

// ValidateAndSanitizeSortOrder validates and sanitizes sort order.
// Only allows ASC or DESC (case-insensitive). Returns "DESC" as default if invalid.
//
// Parameters:
//   - sortOrder: The sort order string from user input
//
// Returns:
//   - string: Validated sort order (ASC, DESC, or default DESC)
func ValidateAndSanitizeSortOrder(sortOrder string) string {
	sortOrderUpper := strings.ToUpper(strings.TrimSpace(sortOrder))
	if sortOrderUpper == "ASC" || sortOrderUpper == "DESC" {
		return sortOrderUpper
	}
	return "DESC" // Default to DESC if invalid
}

// ValidateAndSanitizeSortColumn validates sort column against a whitelist to prevent SQL injection.
// Returns the validated column name or empty string if invalid.
//
// Parameters:
//   - sortBy: The sort column from user input
//   - allowedColumns: Whitelist of allowed column names (e.g., []string{"id", "name", "created_at"})
//   - prefix: Optional table/alias prefix (e.g., "role.", "permission.")
//
// Returns:
//   - string: Validated column name with prefix, or empty string if invalid
//
// Example:
//
//	validated := ValidateAndSanitizeSortColumn("name", []string{"id", "name", "created_at"}, "role.")
//	// Returns: "role.name"
//
//	invalid := ValidateAndSanitizeSortColumn("'; DROP TABLE--", []string{"id", "name"}, "")
//	// Returns: ""
func ValidateAndSanitizeSortColumn(sortBy string, allowedColumns []string, prefix string) string {
	// Remove whitespace and convert to lowercase for comparison
	sortByClean := strings.TrimSpace(strings.ToLower(sortBy))

	// Check if sortBy is in the whitelist
	for _, allowed := range allowedColumns {
		if strings.ToLower(allowed) == sortByClean {
			return prefix + allowed
		}
	}

	// Return empty string if not found in whitelist
	return ""
}

// ValidatePaginationParams validates and sanitizes pagination parameters.
// Ensures PerPage and Page are positive integers within reasonable limits.
//
// Parameters:
//   - page: Page number (1-indexed)
//   - perPage: Items per page
//   - maxPerPage: Maximum allowed items per page (default: 100)
//
// Returns:
//   - validatedPage: Validated page number (minimum 1)
//   - validatedPerPage: Validated per page (minimum 1, maximum maxPerPage)
//
// Example:
//
//	page, perPage := ValidatePaginationParams(0, 200, 100)
//	// Returns: page=1, perPage=100
func ValidatePaginationParams(page, perPage, maxPerPage int) (validatedPage, validatedPerPage int) {
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10 // Default to 10 if invalid
	}

	if maxPerPage > 0 && perPage > maxPerPage {
		perPage = maxPerPage
	}

	return page, perPage
}

// BuildSearchConditionForRawSQL builds a search condition clause and arguments for raw SQL queries.
// This is useful when you need to use ARRAY_AGG or other complex SQL features that require raw queries.
// Uses the same logic as ApplySearchCondition but returns SQL clause string for raw queries.
//
// Parameters:
//   - searchQuery: The search string to match against
//   - searchColumns: List of column names to search in (e.g., []string{"role.name", "pg.module"})
//   - startArgIndex: Starting index for PostgreSQL parameter placeholders (default: 1)
//   - clauseType: Type of clause - "WHERE" or "HAVING" (default: "HAVING" for GROUP BY queries)
//
// Returns:
//   - clause: SQL clause string (e.g., " HAVING (role.name ILIKE $1 OR pg.module ILIKE $2)")
//   - args: Arguments for parameter binding (one per column)
//
// Example:
//
//	clause, args := BuildSearchConditionForRawSQL("john", []string{"role.name", "pg.module"}, 1, "HAVING")
//	// Returns: clause=" HAVING (role.name ILIKE $1 OR pg.module ILIKE $2)", args=["%john%", "%john%"]
func BuildSearchConditionForRawSQL(searchQuery string, searchColumns []string, startArgIndex int, clauseType string) (clause string, args []interface{}) {
	if searchQuery == "" || len(searchColumns) == 0 {
		return "", []interface{}{}
	}

	if clauseType == "" {
		clauseType = "HAVING" // Default for GROUP BY queries
	}

	if startArgIndex < 1 {
		startArgIndex = 1
	}

	// Build OR conditions for each column using PostgreSQL numbered parameters
	// Each column gets its own parameter index (matching ApplySearchCondition behavior)
	conditions := make([]string, len(searchColumns))
	searchPattern := "%" + searchQuery + "%"
	currentArgIndex := startArgIndex

	for i, column := range searchColumns {
		conditions[i] = column + " ILIKE $" + fmt.Sprintf("%d", currentArgIndex)
		args = append(args, searchPattern)
		currentArgIndex++
	}

	// Combine with OR
	whereClause := " " + clauseType + " (" + strings.Join(conditions, " OR ") + ")"

	return whereClause, args
}
