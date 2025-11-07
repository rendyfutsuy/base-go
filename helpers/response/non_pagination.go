package response

type NonPaginationResponse struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (response NonPaginationResponse) SetResponse(data interface{}) (result NonPaginationResponse, err error) {
	result = response

	result.Data = data

	// default value
	result.Status = 200

	result.Message = "page Successfully loaded"

	return result, err
}

// SetErrorResponse creates a NonPaginationResponse for error responses
func SetErrorResponse(statusCode int, message string) NonPaginationResponse {
	return NonPaginationResponse{
		Status:  statusCode,
		Data:    nil,
		Message: message,
	}
}
