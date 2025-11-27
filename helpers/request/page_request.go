package request

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NeedSubqueryPredefine is an interface for repositories that need to define search columns and EXISTS subqueries
type NeedSubqueryPredefine interface {
	GetSearchColumns() []string
	GetSearchExistsSubqueries() []string
}

// NeedFilterPredefine is an interface for repositories that need to apply filters to queries
// The ApplyFilters method should apply all filter conditions from the filter DTO to the query
// The filter parameter is interface{} to allow different filter types per repository
type NeedFilterPredefine interface {
	ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB
}

// PostgreSQL pg_trgm similarity threshold (default). Adjust if needed.
const defaultSimilarityThreshold = 0.45

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

// ApplySearchCondition applies search condition to a GORM query using ILIKE with per-character matching and space-agnostic search.
// The search query is processed character by character. All spaces are removed from both the search query and database columns before comparison.
// This allows "600D" to match "6 00D", "600 D", and "600D".
// If any character in the search query is found in any column (OR between characters, OR between columns), the record will be included in the results.
// If searchQuery is empty, the query is returned unchanged.
//
// Parameters:
//   - query: The GORM query builder
//   - searchQuery: The search string to match against (processed character by character, spaces are ignored)
//   - searchColumns: List of column names to search in (e.g., []string{"usr.full_name", "usr.email"})
func ApplySearchCondition(query *gorm.DB, searchQuery string, searchColumns []string) *gorm.DB {
	if searchQuery == "" || len(searchColumns) == 0 {
		return query
	}

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return query
	}

	// Remove all spaces from search query to enable space-agnostic search
	// and add no space string to words
	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)
	fmt.Println("words: ", words)

	wordClauses := make([]string, 0, len(words))
	args := make([]interface{}, 0)

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := "%" + strings.ToLower(w) + "%"
		columnConditions := make([]string, 0, len(searchColumns))

		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			// SIMILARITY(LOWER(REPLACE(column, ' ', '')), wLower) >= threshold
			columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
			args = append(args, wLower)
		}

		if len(columnConditions) > 0 {
			wordClauses = append(wordClauses, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordClauses) == 0 {
		return query
	}

	whereClause := "(" + strings.Join(wordClauses, " OR ") + ")"
	query = ApplyRelevanceSorting(query, searchQuery, searchColumns)
	return query.Where(whereClause, args...)
}

// ApplySearchConditionWithSubqueries applies search condition to a GORM query using per-character matching (spaces removed).
// Each character must match at least one column or EXISTS subquery (OR between columns/subqueries, OR between characters).
// This function supports both direct column searches and EXISTS subqueries for related tables.
// SearchColumnConfig represents a single search column configuration
// It can be either a simple column name or a custom SQL condition with placeholders
type SearchColumnConfig struct {
	// Column is a simple column name (e.g., "s.supplier_code")
	// If provided, it will be used as: column + " ILIKE ?"
	Column string

	// Condition is a custom SQL condition with placeholders (e.g., "EXISTS (SELECT 1 FROM supplier_contacts sc WHERE sc.supplier_id = s.id AND sc.deleted_at IS NULL AND sc.phone_number ILIKE ?)")
	// If provided, this will be used instead of Column
	// The condition should contain exactly one "?" placeholder for the search pattern
	Condition string
}

// ApplySearchConditionWithSubquery applies search condition to a GORM query with support for both simple columns and subqueries.
// This is an enhanced version of ApplySearchCondition that supports custom SQL conditions like EXISTS subqueries.
// The search query is processed character by character. All spaces are removed from both the search query and database columns before comparison.
// This allows "600D" to match "6 00D", "600 D", and "600D".
// Each character in the search query must be found in at least one column or condition (AND between characters, OR between columns/conditions).
// If all characters match any column or condition (after removing spaces), the record will be included in the results.
// If searchQuery is empty, the query is returned unchanged.
//
// Parameters:
//   - query: The GORM query builder
//   - searchQuery: The search string to match against (words separated by spaces)
//   - searchColumns: List of column names to search in (e.g., []string{"c.customer_code", "c.customer_name"})
//   - existsSubqueries: List of EXISTS subquery templates with placeholder ? for search pattern
//     (e.g., []string{"EXISTS (SELECT 1 FROM customer_contacts cc WHERE cc.customer_id = c.id AND cc.deleted_at IS NULL AND cc.phone_number ILIKE ?)"})
//   - searchQuery: The search string to match against (processed character by character, spaces are ignored)
//   - searchConfigs: List of SearchColumnConfig - each can be a simple column or a custom SQL condition
//
// Returns:
//   - *gorm.DB: The query with search condition applied (if searchQuery is not empty)
func ApplySearchConditionWithSubqueries(query *gorm.DB, searchQuery string, searchColumns []string, existsSubqueries []string) *gorm.DB {
	if searchQuery == "" || (len(searchColumns) == 0 && len(existsSubqueries) == 0) {
		return query
	}

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return query
	}

	// Remove all spaces from search query to enable space-agnostic search
	// and add no space string to words
	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)
	fmt.Println("words: ", words)

	wordClauses := make([]string, 0, len(words))
	args := make([]interface{}, 0)

	reILIKE := regexp.MustCompile(`([a-zA-Z0-9_.]+)\s+ILIKE\s+\?`)

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := "%" + strings.ToLower(w) + "%"
		columnConditions := make([]string, 0, len(searchColumns)+len(existsSubqueries))

		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
			args = append(args, wLower)
		}

		for _, existsSubquery := range existsSubqueries {
			// Replace ILIKE ? with SIMILARITY(LOWER(REPLACE(col,' ', '')), ?) >= threshold
			modified := reILIKE.ReplaceAllString(existsSubquery, "SIMILARITY(LOWER(REPLACE($1, ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
			columnConditions = append(columnConditions, modified)
			args = append(args, wLower)
		}

		if len(columnConditions) > 0 {
			wordClauses = append(wordClauses, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordClauses) == 0 {
		return query
	}

	whereClause := "(" + strings.Join(wordClauses, " OR ") + ")"
	query = ApplyRelevanceSorting(query, searchQuery, searchColumns)
	return query.Where(whereClause, args...)
}

func ApplySearchConditionWithSubquery(query *gorm.DB, searchQuery string, searchConfigs []SearchColumnConfig) *gorm.DB {
	if searchQuery == "" || len(searchConfigs) == 0 {
		return query
	}

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return query
	}

	// Remove all spaces from search query to enable space-agnostic search
	// and add no space string to words
	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)
	fmt.Println("words: ", words)

	wordClauses := make([]string, 0, len(words))
	args := make([]interface{}, 0)

	reILIKE := regexp.MustCompile(`([a-zA-Z0-9_.]+)\s+ILIKE\s+\?`)

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := "%" + strings.ToLower(w) + "%"
		columnConditions := make([]string, 0, len(searchConfigs))

		for _, config := range searchConfigs {
			if config.Condition != "" {
				modifiedCondition := reILIKE.ReplaceAllString(config.Condition, "SIMILARITY(LOWER(REPLACE($1, ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
				columnConditions = append(columnConditions, modifiedCondition)
				args = append(args, wLower)
			} else if config.Column != "" {
				columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
				if !columnRegex.MatchString(config.Column) {
					continue
				}
				columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+config.Column+", ' ', '')), ?) >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
				args = append(args, wLower)
			}
		}

		if len(columnConditions) > 0 {
			wordClauses = append(wordClauses, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordClauses) == 0 {
		return query
	}

	whereClause := "(" + strings.Join(wordClauses, " OR ") + ")"
	cols := make([]string, 0)
	for _, c := range searchConfigs {
		if c.Column != "" {
			cols = append(cols, c.Column)
		}
	}
	query = ApplyRelevanceSorting(query, searchQuery, cols)
	return query.Where(whereClause, args...)
}

func ApplyRelevanceSorting(query *gorm.DB, searchQuery string, searchColumns []string) *gorm.DB {
	if searchQuery == "" || len(searchColumns) == 0 {
		return query
	}
	tokens := strings.Fields(searchQuery)
	collapsed := strings.ReplaceAll(strings.TrimSpace(searchQuery), " ", "")
	relevanceTerms := make([]string, 0)
	relevanceArgs := make([]interface{}, 0)
	addTerms := func(token string) {
		t := strings.TrimSpace(token)
		if t == "" {
			return
		}
		wNorm := strings.ToLower(strings.ReplaceAll(t, " ", ""))
		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			relevanceTerms = append(relevanceTerms, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), ?)")
			relevanceArgs = append(relevanceArgs, wNorm)
		}
	}
	for _, t := range tokens {
		addTerms(t)
	}
	if collapsed != "" {
		addTerms(collapsed)
	}
	if len(relevanceTerms) == 0 {
		return query
	}
	scoreSQL := "(" + strings.Join(relevanceTerms, " + ") + ") DESC"
	return query.Order(clause.Expr{SQL: scoreSQL, Vars: relevanceArgs})
}

// ApplySearchConditionWithSubqueriesFromInterface applies search condition using NeedSubqueryPredefine interface
// This is a convenience wrapper that calls ApplySearchConditionWithSubqueries with values from the interface
func ApplySearchConditionWithSubqueriesFromInterface(query *gorm.DB, searchQuery string, searcher NeedSubqueryPredefine) *gorm.DB {
	if searcher == nil {
		return query
	}
	return ApplySearchConditionWithSubqueries(query, searchQuery, searcher.GetSearchColumns(), searcher.GetSearchExistsSubqueries())
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

// BuildNaturalSortExpression builds a natural/numeric sorting expression for PostgreSQL.
// This is useful for sorting columns that contain numbers mixed with text (e.g., "String 1", "String 10").
// The function sorts by text part (prefix before first number) first alphabetically, then by first numeric part numerically, then by original column.
// This ensures proper natural sorting where:
// - Numbers are sorted numerically (1, 2, 10) not lexicographically (1, 10, 2)
// - Strings with the same prefix and number are sorted by original column (e.g., "String 6" before "String 6 UPDATED")
// - Strings without numbers appear after strings with numbers for the same text prefix (for ASC)
//
// Parameters:
//   - column: The column name to sort (e.g., "role.name")
//   - sortOrder: Sort order ("ASC" or "DESC")
//   - useNaturalSort: Whether to use natural sorting (if false, returns standard sorting expression)
//
// Returns:
//   - string: SQL sorting expression
func BuildNaturalSortExpression(column, sortOrder string, useNaturalSort bool) string {
	// Validate column name to prevent SQL injection (only allow alphanumeric, underscore, and dot)
	// This is a defense-in-depth measure - column names should already be validated by ValidateAndSanitizeSortColumn
	columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !columnRegex.MatchString(column) {
		// If column name is invalid, return safe default
		return "created_at " + ValidateAndSanitizeSortOrder(sortOrder)
	}

	// Validate sort order
	sortOrder = ValidateAndSanitizeSortOrder(sortOrder)
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	if !useNaturalSort {
		return column + " " + sortOrder
	}

	// Use natural sorting: sort by text part (prefix before first number) first, then by numeric part, then by original column
	// This ensures proper natural sorting where:
	// - Strings are grouped by text part (prefix before first number)
	// - For the same text part, strings are sorted by first numeric value found
	// - For the same numeric value, strings are sorted by original column (to handle cases like "String 6" vs "String 6 UPDATED")
	// Extract text part before first number (prefix text before any number appears)
	textPart := fmt.Sprintf("TRIM(SUBSTRING(%s FROM '^([^0-9]*)'))", column)
	// Extract first numeric part (extract first sequence of digits, convert to bigint, NULL if no digits)
	// This gets the first number found in the string
	numericPart := fmt.Sprintf("(CASE WHEN (regexp_match(%s, '[0-9]+'))[1] IS NULL THEN NULL ELSE ((regexp_match(%s, '[0-9]+'))[1])::bigint END)", column, column)

	// Sort by text part first (alphabetically), then by numeric part, then by original column
	// For ASC: text part ASC, then numeric part ASC NULLS LAST (strings without numbers appear after strings with numbers for the same text part), then column ASC
	// For DESC: text part DESC, then numeric part DESC NULLS LAST (strings without numbers appear after strings with numbers for the same text part), then column DESC
	if strings.EqualFold(sortOrder, "ASC") {
		return fmt.Sprintf("%s ASC, %s %s NULLS LAST, %s ASC", textPart, numericPart, sortOrder, column)
	} else {
		return fmt.Sprintf("%s DESC, %s %s NULLS LAST, %s DESC", textPart, numericPart, sortOrder, column)
	}
}

// determineSortExpression is an internal helper that determines the sort expression based on request and config.
// This logic is shared between ApplyPagination and BuildSortExpressionForRawSQL.
func determineSortExpression(req PageRequest, config PaginationConfig, customSortMapping func(string) string) string {
	// Set defaults
	if config.DefaultSortBy == "" {
		config.DefaultSortBy = "created_at"
	}
	if config.DefaultSortOrder == "" {
		config.DefaultSortOrder = "DESC"
	}

	// Determine sort column
	sortBy := config.DefaultSortBy
	if req.SortBy != "" {
		// Use custom mapping if provided, otherwise use config.SortMapping
		mappingFunc := customSortMapping
		if mappingFunc == nil {
			mappingFunc = config.SortMapping
		}

		if mappingFunc != nil {
			// Use custom mapping function if provided
			mapped := mappingFunc(req.SortBy)
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

	// Check if natural sorting is needed for this column
	useNaturalSort := false
	if len(config.NaturalSortColumns) > 0 {
		for _, naturalCol := range config.NaturalSortColumns {
			if naturalCol == sortBy {
				useNaturalSort = true
				break
			}
		}
	}

	return BuildNaturalSortExpression(sortBy, sortOrder, useNaturalSort)
}

// BuildSortExpressionForRawSQL builds a sort expression for raw SQL queries.
// This is useful for raw SQL queries where you need to build the ORDER BY clause manually.
// Uses the same logic as ApplyPagination.
//
// Parameters:
//   - req: PageRequest containing sort_by and sort_order
//   - config: PaginationConfig with default values and validation rules
//   - customSortMapping: Optional custom sort column mapping function (if nil, uses config.SortMapping)
//
// Returns:
//   - string: The complete sort expression (with natural sorting if applicable)
func BuildSortExpressionForRawSQL(req PageRequest, config PaginationConfig, customSortMapping func(string) string) string {
	return determineSortExpression(req, config, customSortMapping)
}

// BuildSortExpressionForExport builds a sort expression for export queries.
// This is similar to BuildSortExpressionForRawSQL but designed for export queries that don't use PageRequest.
//
// Parameters:
//   - sortBy: The sort column name (e.g., "employee_name" or "expedition_name")
//   - sortOrder: Sort order ("ASC" or "DESC")
//   - defaultSortBy: Default sort column if sortBy is empty
//   - defaultSortOrder: Default sort order if sortOrder is empty
//   - sortMapping: Function to map sortBy to actual database column (e.g., mapEmployeeIndexSortColumn)
//   - naturalSortColumns: List of columns that should use natural sorting (e.g., []string{"e.employee_name", "e.expedition_name"})
//
// Returns:
//   - string: The complete sort expression (with natural sorting if applicable)
func BuildSortExpressionForExport(
	sortBy, sortOrder string,
	defaultSortBy, defaultSortOrder string,
	sortMapping func(string) string,
	naturalSortColumns []string,
) string {
	// Use default if sortBy is empty
	if sortBy == "" {
		sortBy = defaultSortBy
	}

	// Map sortBy to actual column
	mappedSortBy := sortBy
	if sortMapping != nil {
		if mapped := sortMapping(sortBy); mapped != "" {
			mappedSortBy = mapped
		}
	}

	// Use default if sortOrder is empty
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}
	sortOrder = ValidateAndSanitizeSortOrder(sortOrder)
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}

	// Check if natural sorting is needed for this column
	useNaturalSort := false
	if len(naturalSortColumns) > 0 {
		for _, naturalCol := range naturalSortColumns {
			if naturalCol == mappedSortBy {
				useNaturalSort = true
				break
			}
		}
	}

	return BuildNaturalSortExpression(mappedSortBy, sortOrder, useNaturalSort)
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
// The search query is processed character by character. All spaces are removed from both the search query and database columns before comparison.
// This allows "600D" to match "6 00D", "600 D", and "600D".
// If any character in the search query is found in any column (OR between characters, OR between columns), the record will be included in the results.
//
// Parameters:
//   - searchQuery: The search string to match against (processed character by character, spaces are ignored)
//   - searchColumns: List of column names to search in (e.g., []string{"role.name", "pg.module"})
//   - startArgIndex: Starting index for PostgreSQL parameter placeholders (default: 1)
//   - clauseType: Type of clause - "WHERE" or "HAVING" (default: "HAVING" for GROUP BY queries)
//
// Returns:
//   - clause: SQL clause string (e.g., " HAVING ((REPLACE(role.name, ' ', ”) ILIKE $1 OR REPLACE(pg.module, ' ', ”) ILIKE $2) AND ...)")
//   - args: Arguments for parameter binding (one per column per character)
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

	words := strings.Fields(searchQuery)
	if len(words) == 0 {
		return "", []interface{}{}
	}

	// Remove all spaces from search query to enable space-agnostic search
	// and add no space string to words
	searchQueryNoSpaces := strings.ReplaceAll(searchQuery, " ", "")
	words = append(words, searchQueryNoSpaces)
	fmt.Println("words: ", words)

	wordConditions := make([]string, 0, len(words))
	currentArgIndex := startArgIndex

	for _, w := range words {
		w = strings.TrimSpace(w)
		if w == "" {
			continue
		}
		wLower := strings.ToLower(w)
		columnConditions := make([]string, 0, len(searchColumns))

		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			// SIMILARITY(LOWER(REPLACE(column, ' ', '')), $N) >= threshold
			columnConditions = append(columnConditions, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), $"+fmt.Sprintf("%d", currentArgIndex)+") >= "+fmt.Sprintf("%0.2f", defaultSimilarityThreshold))
			args = append(args, wLower)
			currentArgIndex++
		}

		if len(columnConditions) > 0 {
			wordConditions = append(wordConditions, "("+strings.Join(columnConditions, " OR ")+")")
		}
	}

	if len(wordConditions) == 0 {
		return "", []interface{}{}
	}

	clause = " " + clauseType + " (" + strings.Join(wordConditions, " AND ") + ")"
	return clause, args
}

func BuildRelevanceOrderByForRawSQL(searchQuery string, searchColumns []string, startArgIndex int) (orderBy string, args []interface{}) {
	if searchQuery == "" || len(searchColumns) == 0 {
		return "", []interface{}{}
	}
	tokens := strings.Fields(searchQuery)
	terms := make([]string, 0)
	args = make([]interface{}, 0)
	current := startArgIndex
	for _, t := range tokens {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		wNorm := strings.ToLower(strings.ReplaceAll(t, " ", ""))
		for _, column := range searchColumns {
			columnRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
			if !columnRegex.MatchString(column) {
				continue
			}
			terms = append(terms, "SIMILARITY(LOWER(REPLACE("+column+", ' ', '')), $"+fmt.Sprintf("%d", current)+")")
			args = append(args, wNorm)
			current++
		}
	}
	if len(terms) == 0 {
		return "", []interface{}{}
	}
	orderBy = " ORDER BY (" + strings.Join(terms, " + ") + ") DESC"
	return orderBy, args
}

// PaginationConfig holds configuration for pagination
type PaginationConfig struct {
	DefaultSortBy      string              // Default sort column (e.g., "usr.created_at")
	DefaultSortOrder   string              // Default sort order (e.g., "DESC")
	AllowedColumns     []string            // Allowed sort columns for validation (e.g., []string{"id", "name", "created_at"})
	ColumnPrefix       string              // Table/alias prefix (e.g., "usr.", "role.")
	MaxPerPage         int                 // Maximum items per page (default: 100)
	SortMapping        func(string) string // Optional custom sort column mapping function
	NaturalSortColumns []string            // Columns that should use natural/numeric sorting (e.g., []string{"role.name"})
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

	// Determine sort expression using shared logic
	sortExpression := determineSortExpression(req, config, nil)

	// Apply sorting and pagination
	// Use Scan() as it works for both standard GORM queries and custom SELECT queries with joins
	err = query.
		Order(sortExpression).
		Limit(validatedPerPage).
		Offset(offset).
		Scan(result).Error

	return total, nil
}
