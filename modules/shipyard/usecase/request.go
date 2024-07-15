package usecase

import (
	"fmt"
	"strings"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
)

const (
	fetchQueryFilter = "deleted_at IS NULL"
)

var (
	fetchQueryFormatWithFilter = fetchQueryFilter + " AND (%s)"
	fetchQuerySearch           = "CONCAT_WS(' ', %s) ILIKE '%%' || '%s' || '%%'"
)

func constructSearchQuery(searchText string, fields ...string) (result string) {
	joinFields := ""
	if len(fields) == 0 {
		joinFields = "code"
	} else {
		joinFields = strings.Join(fields, ",")
	}
	result = fmt.Sprintf(fetchQuerySearch, joinFields, searchText)
	return
}

// fetchRepoRequest is a struct that implements the QueryRequest interface.
type fetchRepoRequest struct {
	Condition string // Condition is the condition for the query request.
	Params    []any  // Params is the parameters for the query request.
	Paginate  bool   // IsPaginate is the flag to decide if the fetch is need pagination or no.
	PageSize  int    // PageSize is the page size if needed.
	Page      int    // Page is page number.
	SortField string // SortField is field to sort query
	SortOrder string // SortOrder is order to sort
}

// GetCondition is a method of the fetchRepoRequest struct.
// It returns the condition of the query request.
func (req *fetchRepoRequest) GetCondition() (result string) {
	return req.Condition
}

// GetParam is a method of the fetchRepoRequest struct.
// It returns the parameters of the query request.
func (req *fetchRepoRequest) GetParam() (result []any) {
	if len(req.Params) == 0 {
		return nil
	}

	return req.Params
}

// IsPaginate is a method that returns pagination settings.
func (req *fetchRepoRequest) IsPaginate() (isPaginate bool, limit, offset int) {
	isPaginate = req.Paginate
	limit = req.PageSize
	offset = (req.Page - 1) * req.PageSize
	return
}

// GetSort is a function that gives us the sorting details.
// 'field' tells us what we're sorting by.
// 'order' tells us the direction of sorting (like smallest to biggest).
func (req *fetchRepoRequest) GetSort() (field, order string) {
	field = req.SortField // Gets the sorting field from the request
	order = req.SortOrder // Gets the sorting order from the request
	return                // Gives back the sorting field and order
}

// fetchRepoRequesFromPageRequest is a function that takes a page request and turns it into a repository request.
// It sets up the conditions for fetching data from a repository.
func fetchRepoRequesFromPageRequest(req request.PageRequest) (result fetchRepoRequest) {
	result.Condition = fetchQueryFilter // Sets the condition for fetching data
	if req.Search != "" {
		result.Condition = fmt.Sprintf(fetchQueryFormatWithFilter, constructSearchQuery(req.Search, "name", "code", "yard"))
	}
	if req.PerPage > 0 { // Checks if the number of items per page is more than 0
		result.Paginate = true        // If so, it turns on pagination
		result.Page = req.Page        // Sets the page number
		result.PageSize = req.PerPage // Sets the number of items per page
		result.SortField = req.SortBy // Sets the field to sort by
		if result.SortField == "" {   // If no sort field is provided, it defaults to "created_at"
			result.SortField = "created_at"
		}
		result.SortOrder = req.SortOrder // Sets the order of sorting
		if result.SortOrder == "" {      // If no sort order is provided, it defaults to "ASC" (ascending)
			result.SortOrder = "ASC"
		}
	}
	return // Returns the repository request
}
