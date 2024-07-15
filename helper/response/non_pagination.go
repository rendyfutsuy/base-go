package response

type NonPaginationResponse struct {
	Data interface{}    `json:"data"`
}

func (response NonPaginationResponse) SetResponse(data interface{}) (result NonPaginationResponse, err error) {
	result = response

	result.Data = data

	return result, err
}
