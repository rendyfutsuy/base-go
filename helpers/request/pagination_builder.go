package request

import (
	"fmt"
	"strings"
)

type PaginationConfig struct {
	DefaultSortBy      string
	DefaultSortOrder   string
	AllowedColumns     []string
	ColumnPrefix       string
	MaxPerPage         int
	SortMapping        func(string) string
	NaturalSortColumns []string
}

// ValidateAndSanitizeSortOrder validates the sort order and returns a sanitized version of the input string.
// It trims whitespace from the input string, converts it to uppercase, and checks if it is either "ASC" or "DESC".
// If the input is valid, it returns the sanitized string. Otherwise, it returns "DESC".
func ValidateAndSanitizeSortOrder(sortOrder string) string {
	sortOrderUpper := strings.ToUpper(strings.TrimSpace(sortOrder))
	if sortOrderUpper == "ASC" || sortOrderUpper == "DESC" {
		return sortOrderUpper
	}
	return "DESC"
}

// ValidateAndSanitizeSortColumn validates the sort column and returns a sanitized version of the input string.
// It trims whitespace from the input string, converts it to lowercase, and checks if it is in the allowedColumns.
// If the input is valid, it returns the sanitized string prefixed with the given prefix. Otherwise, it returns an empty string.
func ValidateAndSanitizeSortColumn(sortBy string, allowedColumns []string, prefix string) string {
	sortByClean := strings.TrimSpace(strings.ToLower(sortBy))
	for _, allowed := range allowedColumns {
		if strings.ToLower(allowed) == sortByClean {
			return prefix + allowed
		}
	}
	return ""
}

// BuildNaturalSortExpression builds a natural sort expression based on the given column, sort order, and flag.
// If the column name does not match the regex for a valid column name, it returns a default sort expression for the "created_at" column.
// The function takes into account the sort order and whether to use natural sort or not.
// If natural sort is enabled, the function returns a sort expression that first sorts by the text part of the column (everything before the first digit), then by the numeric part of the column (the first sequence of digits), and finally by the column itself.
// If natural sort is disabled, the function returns a simple sort expression for the column with the given sort order.
// The function also handles the case where the sort order is not provided, in which case it defaults to "DESC".
func BuildNaturalSortExpression(column, sortOrder string, useNaturalSort bool) string {
	// 2026/01/27: comment this logic to disable sort without aggregated column.
	if !ColumnNameRegex.MatchString(column) {
		return "created_at " + ValidateAndSanitizeSortOrder(sortOrder)
	}
	sortOrder = ValidateAndSanitizeSortOrder(sortOrder)
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	if !useNaturalSort {
		return column + " " + sortOrder
	}
	textPart := fmt.Sprintf("TRIM(SUBSTRING(%s FROM '^([^0-9]*)'))", column)
	numericPart := fmt.Sprintf("(CASE WHEN (regexp_match(%s, '[0-9]+'))[1] IS NULL THEN NULL ELSE ((regexp_match(%s, '[0-9]+'))[1])::bigint END)", column, column)
	if strings.EqualFold(sortOrder, "ASC") {
		return fmt.Sprintf("%s ASC, %s %s NULLS LAST, %s ASC", textPart, numericPart, sortOrder, column)
	} else {
		return fmt.Sprintf("%s DESC, %s %s NULLS LAST, %s DESC", textPart, numericPart, sortOrder, column)
	}
}

// determineSortExpression determines the sort expression for a given page request and pagination config.
// It uses the provided custom sort mapping function if not nil, otherwise it uses the pagination config's sort mapping.
// If the sort by is empty, it uses the pagination config's default sort by.
// If the sort order is empty, it uses the pagination config's default sort order.
// If the sort by is in the pagination config's natural sort columns, it uses natural sort instead of default sort.
func determineSortExpression(req PageRequest, config PaginationConfig, customSortMapping func(string) string) string {
	if config.DefaultSortBy == "" {
		config.DefaultSortBy = "created_at"
	}
	if config.DefaultSortOrder == "" {
		config.DefaultSortOrder = "DESC"
	}
	sortBy := config.DefaultSortBy
	if req.SortBy != "" {
		mappingFunc := customSortMapping
		if mappingFunc == nil {
			mappingFunc = config.SortMapping
		}
		if mappingFunc != nil {
			mapped := mappingFunc(req.SortBy)
			if mapped != "" {
				sortBy = mapped
			}
		} else if len(config.AllowedColumns) > 0 {
			validated := ValidateAndSanitizeSortColumn(req.SortBy, config.AllowedColumns, config.ColumnPrefix)
			if validated != "" {
				sortBy = validated
			}
		}
	}
	sortOrder := config.DefaultSortOrder
	if req.SortOrder != "" {
		validated := ValidateAndSanitizeSortOrder(req.SortOrder)
		if validated != "" {
			sortOrder = validated
		}
	}
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

// BuildSortExpressionForRawSQL returns the sort expression for raw SQL queries based on the page request and pagination config.
// It uses the provided custom sort mapping function if not nil, otherwise it uses the pagination config's sort mapping.
// If the sort by is empty, it uses the pagination config's default sort by.
// If the sort order is empty, it uses the pagination config's default sort order.
// If the sort by is in the pagination config's natural sort columns, it uses natural sort instead of default sort.
func BuildSortExpressionForRawSQL(req PageRequest, config PaginationConfig, customSortMapping func(string) string) string {
	return determineSortExpression(req, config, customSortMapping)
}

// BuildSortExpressionForExport returns the sort expression for export queries based on the provided sort by and sort order.
// It uses the provided custom sort mapping function if not nil, otherwise it uses the pagination config's sort mapping.
// If the sort by is empty, it uses the pagination config's default sort by.
// If the sort order is empty, it uses the pagination config's default sort order.
// If the sort by is in the pagination config's natural sort columns, it uses natural sort instead of default sort.
func BuildSortExpressionForExport(sortBy, sortOrder string, defaultSortBy, defaultSortOrder string, sortMapping func(string) string, naturalSortColumns []string) string {
	if sortBy == "" {
		sortBy = defaultSortBy
	}
	mappedSortBy := sortBy
	if sortMapping != nil {
		mapped := sortMapping(sortBy)
		if mapped != "" {
			mappedSortBy = mapped
		} else {
			mappedSortBy = defaultSortBy
		}
	}
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}
	sortOrder = ValidateAndSanitizeSortOrder(sortOrder)
	if sortOrder == "" {
		sortOrder = defaultSortOrder
	}
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

// ValidatePaginationParams validates the pagination parameters (page and perPage) and returns the validated values.
// If page is less than 1, it sets page to 1.
// If perPage is less than 1, it sets perPage to 10.
// If perPage is greater than maxPerPage and maxPerPage is greater than 0, it sets perPage to maxPerPage.
func ValidatePaginationParams(page, perPage, maxPerPage int) (validatedPage, validatedPerPage int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if maxPerPage > 0 && perPage > maxPerPage {
		perPage = maxPerPage
	}
	return page, perPage
}
