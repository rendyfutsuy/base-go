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

// ApplySearchCondition applies search condition to a GORM query using ILIKE with per-word matching.
// Each word in the search query (separated by spaces) is searched independently (OR between words, OR between columns).
// If any word matches any column, the record will be included in the results.
// If searchQuery is empty, the query is returned unchanged.
//
// Parameters:
//   - query: The GORM query builder
//   - searchQuery: The search string to match against (words separated by spaces)
//   - searchColumns: List of column names to search in (e.g., []string{"usr.full_name", "usr.email"})
//
// Returns:
//   - *gorm.DB: The query with search condition applied (if searchQuery is not empty)
//
// Example:
//
//	query = ApplySearchCondition(query, "Add Role", []string{"usr.full_name", "usr.email", "rl.name"})
//	// Results in: WHERE ((usr.full_name ILIKE '%Add%' OR usr.email ILIKE '%Add%' OR rl.name ILIKE '%Add%')
//	//            OR (usr.full_name ILIKE '%Role%' OR usr.email ILIKE '%Role%' OR rl.name ILIKE '%Role%'))
func ApplySearchCondition(query *gorm.DB, searchQuery string, searchColumns []string) *gorm.DB {
	if searchQuery == "" || len(searchColumns) == 0 {
		return query
	}

	// Split search query by spaces and filter out empty strings
	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return query
	}

	// Build conditions for each word
	wordConditions := make([]string, 0, len(words))
	args := make([]interface{}, 0)

	for _, word := range words {
		if word == "" {
			continue
		}

		// For each word, build OR conditions for all columns
		columnConditions := make([]string, len(searchColumns))
		searchPattern := "%" + word + "%"

		for i, column := range searchColumns {
			columnConditions[i] = column + " ILIKE ?"
			args = append(args, searchPattern)
		}

		// Combine column conditions with OR and wrap in parentheses
		wordConditions = append(wordConditions, "("+strings.Join(columnConditions, " OR ")+")")
	}

	if len(wordConditions) == 0 {
		return query
	}

	// Combine all word conditions with OR (not AND)
	whereClause := "(" + strings.Join(wordConditions, " OR ") + ")"

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
// Each word in the search query (separated by spaces) is searched independently (OR between words, OR between columns).
// If any word matches any column, the record will be included in the results.
//
// Parameters:
//   - searchQuery: The search string to match against (words separated by spaces)
//   - searchColumns: List of column names to search in (e.g., []string{"role.name", "pg.module"})
//   - startArgIndex: Starting index for PostgreSQL parameter placeholders (default: 1)
//   - clauseType: Type of clause - "WHERE" or "HAVING" (default: "HAVING" for GROUP BY queries)
//
// Returns:
//   - clause: SQL clause string (e.g., " HAVING ((role.name ILIKE $1 OR pg.module ILIKE $2) OR (role.name ILIKE $3 OR pg.module ILIKE $4))")
//   - args: Arguments for parameter binding (one per column per word)
//
// Example:
//
//	clause, args := BuildSearchConditionForRawSQL("Add Role", []string{"role.name", "pg.module"}, 1, "HAVING")
//	// Returns: clause=" HAVING ((role.name ILIKE $1 OR pg.module ILIKE $2) OR (role.name ILIKE $3 OR pg.module ILIKE $4))"
//	//          args=["%Add%", "%Add%", "%Role%", "%Role%"]
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

	// Split search query by spaces and filter out empty strings
	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return "", []interface{}{}
	}

	// Build conditions for each word
	wordConditions := make([]string, 0, len(words))
	currentArgIndex := startArgIndex

	for _, word := range words {
		if word == "" {
			continue
		}

		// For each word, build OR conditions for all columns
		columnConditions := make([]string, len(searchColumns))
		searchPattern := "%" + word + "%"

		for i, column := range searchColumns {
			columnConditions[i] = column + " ILIKE $" + fmt.Sprintf("%d", currentArgIndex)
			args = append(args, searchPattern)
			currentArgIndex++
		}

		// Combine column conditions with OR and wrap in parentheses
		wordConditions = append(wordConditions, "("+strings.Join(columnConditions, " OR ")+")")
	}

	if len(wordConditions) == 0 {
		return "", []interface{}{}
	}

	// Combine all word conditions with OR (not AND)
	whereClause := " " + clauseType + " (" + strings.Join(wordConditions, " OR ") + ")"

	return whereClause, args
}

// PaginationConfig holds configuration for pagination
type PaginationConfig struct {
	DefaultSortBy    string              // Default sort column (e.g., "usr.created_at")
	DefaultSortOrder string              // Default sort order (e.g., "DESC")
	AllowedColumns   []string            // Allowed sort columns for validation (e.g., []string{"id", "name", "created_at"})
	ColumnPrefix     string              // Table/alias prefix (e.g., "usr.", "role.")
	MaxPerPage       int                 // Maximum items per page (default: 100)
	SortMapping      func(string) string // Optional custom sort column mapping function
}

// ApplyPagination applies pagination, sorting, and count logic to a GORM query.
// This is a generic function that handles common pagination logic for index/list endpoints.
//
// Parameters:
//   - query: The GORM query builder (with filters already applied, but before count/pagination)
//   - req: PageRequest containing page, per_page, sort_by, and sort_order
//   - config: PaginationConfig with default values and validation rules
//   - result: Pointer to slice where results will be scanned (e.g., &[]models.User)
//
// Returns:
//   - total: Total count of records (before pagination)
//   - err: Error if pagination fails
//
// Example:
//
//	config := PaginationConfig{
//		DefaultSortBy:    "usr.created_at",
//		DefaultSortOrder: "DESC",
//		AllowedColumns:   []string{"id", "full_name", "email", "created_at"},
//		ColumnPrefix:     "usr.",
//		MaxPerPage:       100,
//	}
//	total, err := ApplyPagination(query, req, config, &users)
func ApplyPagination(query *gorm.DB, req PageRequest, config PaginationConfig, result interface{}) (total int, err error) {
	// Set defaults
	if config.MaxPerPage <= 0 {
		config.MaxPerPage = 100
	}
	if config.DefaultSortBy == "" {
		config.DefaultSortBy = "created_at"
	}
	if config.DefaultSortOrder == "" {
		config.DefaultSortOrder = "DESC"
	}

	// Validate and sanitize pagination parameters
	validatedPage, validatedPerPage := ValidatePaginationParams(req.Page, req.PerPage, config.MaxPerPage)
	offset := (validatedPage - 1) * validatedPerPage

	// Count total (before pagination)
	countQuery := query
	var totalCount int64
	err = countQuery.Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	total = int(totalCount)

	// Determine sort column
	sortBy := config.DefaultSortBy
	if req.SortBy != "" {
		if config.SortMapping != nil {
			// Use custom mapping function if provided
			mapped := config.SortMapping(req.SortBy)
			if mapped != "" {
				sortBy = mapped
			}
		} else if len(config.AllowedColumns) > 0 {
			// Use ValidateAndSanitizeSortColumn if allowedColumns provided
			validated := ValidateAndSanitizeSortColumn(req.SortBy, config.AllowedColumns, config.ColumnPrefix)
			if validated != "" {
				sortBy = validated
			}
		}
	}

	// Determine sort order
	sortOrder := config.DefaultSortOrder
	if req.SortOrder != "" {
		validated := ValidateAndSanitizeSortOrder(req.SortOrder)
		if validated != "" {
			sortOrder = validated
		}
	}

	// Apply sorting and pagination
	// Use Scan() as it works for both standard GORM queries and custom SELECT queries with joins
	err = query.
		Order(sortBy + " " + sortOrder).
		Limit(validatedPerPage).
		Offset(offset).
		Scan(result).Error

	return total, nil
}
