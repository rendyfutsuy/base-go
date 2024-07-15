package request

type PageRequest struct {
	Page      int      `json:"page"`
	PerPage   int      `json:"per_page"`
	Filters   []Filter `json:"filter"`
	Search    string   `json:"search"`
	SortBy    string   `json:"sort_by"`
	SortOrder string   `json:"sort_order"`
}

func NewPageRequest(page, perPage int, search, sortBy, sortOrder string, filters []Filter) *PageRequest {
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
