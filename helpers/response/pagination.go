package response

type PaginationResponse struct {
	Meta PaginationMeta `json:"meta"`
	Data interface{}    `json:"data"`
}

func (response PaginationResponse) SetResponse(data interface{}, dataTotal int, perPage int, currentPage int) (result PaginationResponse, err error) {
	result = response

	result.Data = data

	result.Meta, err = result.Meta.SetMeta(dataTotal, perPage, currentPage)

	return result, err
}
