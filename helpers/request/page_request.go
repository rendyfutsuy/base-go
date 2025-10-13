package request

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iman231/go-help"
)

type PageRequest struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
	Filters    []Filter `json:"filter"`
	Search     string   `json:"search"`
	SortBy     string   `json:"sort_by"`
	SortOrder  string   `json:"sort_order"`
	Projection []string `json:"projections"`
}

func NewPageRequest(page, perPage int, search, sortBy, sortOrder string, filters []Filter, Projections []string) *PageRequest {
	return &PageRequest{
		Page:      page,
		PerPage:   perPage,
		Search:    search,
		Filters:   filters,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

type Filter struct {
	Option    string `json:"option"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
	ValueType string `json:"type"`
}

type Search struct {
	Options  []string `json:"options"`
	Value    string   `json:"value"`
	Operator string   `json:"operator"`
}

func BuildDynamicFilterAndSearch(baseQuery string, filters []Filter, searchs []*Search, isDeleted bool) (string, []interface{}) {
	var whereClauses []string
	var orClauses [][]string
	var args []interface{}

	// Add the condition to check for non-deleted records
	if !isDeleted {
		whereClauses = append(whereClauses, "fac.deleted_at IS NULL")
	} else {
		whereClauses = append(whereClauses, "fac.deleted_at IS NOT NULL")
	}

	for _, condition := range filters {
		// Add the condition to the WHERE clause
		whereClauses = append(whereClauses, fmt.Sprintf("%s %s $%d", condition.Option, condition.Operator, len(args)+1))
		args = append(args, condition.Value)
	}

	availableOperator := []string{
		"LIKE",
		"ILIKE",
		"=",
	}

	if searchs != nil && len(searchs) > 0 {
		for _, search := range searchs {
			var orClause []string
			for _, col := range search.Options {
				operator := "LIKE"

				if search.Operator != "" && help.InArray(search.Operator, availableOperator) {
					operator = search.Operator
				}

				orClause = append(orClause, fmt.Sprintf("%s %s $%d", col, operator, len(args)+1))
				args = append(args, search.Value)
			}
			orClauses = append(orClauses, orClause)
		}
	}

	// If there are conditions, add them to the base query
	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if len(orClauses) > 0 {
		if len(whereClauses) > 0 {
			baseQuery += " AND"
		} else {
			baseQuery += " WHERE"
		}

		for i, clause := range orClauses {
			if i > 0 {
				baseQuery += " AND"
			}

			baseQuery += "(" + strings.Join(clause, " OR ") + ")"
		}

	}

	return baseQuery, args
}

func SetFilter(jsonFilters []string) (filters []Filter, err error) {

	if len(jsonFilters) > 0 {
		for _, jsonData := range jsonFilters {
			if jsonData != "" {
				var filter Filter
				err = json.Unmarshal([]byte(jsonData), &filter)

				if err != nil {
					return
				}

				if filter.Operator == "LIKE" {
					filter.Value = fmt.Sprintf("%%%s%%", filter.Value)
				}

				filters = append(filters, filter)
			}
		}

		return
	}

	return
}
