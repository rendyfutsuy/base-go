package request

import (
	"regexp"

	"gorm.io/gorm"
)

var (
	ColumnNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	ILIKEPatternRegex = regexp.MustCompile(`([a-zA-Z0-9_.]+(?:::[a-zA-Z0-9_]+)?)\s+ILIKE\s+\?`)
)

type NeedFilterPredefine interface {
	ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB
}

type NeedSearchPredefine interface {
	GetSearchColumns() []string
	GetSearchExistsSubqueries() []string
	SimilarityThreshold() *float64 // if this logic somehow return nil, then defaultSimilarityThreshold will be used
}

// PostgreSQL pg_trgm similarity threshold (default). Adjust if needed.
const DefaultSimilarityThreshold = 0.33

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

type SearchColumnConfig struct {
	Column    string
	Condition string
}

// ApplyPagination applies pagination to the query.
// It takes the query, PageRequest, PaginationConfig and result interface as parameters.
// It validates the PageRequest and PaginationConfig, and applies the limit and offset to the query.
// It returns the total count and error.
// The PaginationConfig is used to determine the default sort by and sort order if not provided in the PageRequest.
// The result interface is used to scan the result of the query.
// The function returns the total count and error.
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

// ApplySearchConditionFromInterface applies search condition from the searchQuery to the query.
// It takes the query, searchQuery, searcher and return the modified query.
// The searcher is used to determine the search columns and subqueries.
// The threshold is used to determine the similarity threshold for the search condition.
// If the searcher is nil, the function returns the query unchanged.
func ApplySearchConditionFromInterface(query *gorm.DB, searchQuery string, searcher NeedSearchPredefine) *gorm.DB {
	if searcher == nil {
		return query
	}
	threshold := DefaultSimilarityThreshold
	if t := searcher.SimilarityThreshold(); t != nil {
		threshold = *t
	}
	return ApplySearchConditionWithSubqueriesThreshold(query, searchQuery, searcher.GetSearchColumns(), searcher.GetSearchExistsSubqueries(), threshold)
}
