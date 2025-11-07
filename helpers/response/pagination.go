package response

type PaginationResponse struct {
	Status  int            `json:"status"`
	Meta    PaginationMeta `json:"meta"`
	Data    interface{}    `json:"data"`
	Message string         `json:"message"`
}

func (response PaginationResponse) SetResponse(data interface{}, dataTotal int, perPage int, currentPage int) (result PaginationResponse, err error) {
	result = response

	result.Data = data

	result.Meta, err = result.Meta.SetMeta(dataTotal, perPage, currentPage)

	result.Status = 200

	result.Message = "page Successfully loaded"

	return result, err
}
