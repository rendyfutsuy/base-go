package response

import "fmt"

type PaginationMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
}

func (meta PaginationMeta) SetMeta(total int, perPage int, currentPage int) (result PaginationMeta, err error) {
	if perPage == 0 {
		return result, fmt.Errorf("per_page must be greater than 0")
	}

	lastPage := (total + perPage - 1) / perPage
	from := (currentPage-1)*perPage + 1

	to := from + perPage - 1

	if lastPage == currentPage {
		to = total
	}

	if total == 0 {
		from = 0
		to = 0
	}

	result.Total = total
	result.PerPage = perPage
	result.CurrentPage = currentPage
	result.LastPage = lastPage
	result.From = from
	result.To = to

	return result, err
}
