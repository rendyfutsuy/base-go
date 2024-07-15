package http

import (
	"net/http"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helpers/middleware"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Declare a global variable handler of type *ShipyardHandler.
var handler *ShipyardHandler

// ShipyardHandler is a struct that contains the use case for shipyard,
// a validator for validating data, and an authentication middleware.
type ShipyardHandler struct {
	ShipyardUsecase shipyard.Usecase                   // ShipyardUsecase is the use case for shipyard.
	validator       *validator.Validate                // validator is used for validating data.
	middlewareAuth  middleware.IMiddlewareAuth         // middlewareAuth is the authentication middleware.
	mwPageRequest   _reqContext.IMiddlewarePageRequest // mwPageRequest is a middleware for handling page requests. It processes the request data related to pagination.
}

// GeneralResponse is a struct that represents a general response.
// It contains a message.
type GeneralResponse struct {
	Message string `json:"message"` // Message is the message of the response.
}

// ReturnDataResponse is a struct that represents a response with returned data.
// It contains a status, a meta of type GeneralResponse, and data of any type.
type ReturnDataResponse struct {
	Status      string          `json:"status"`                 // Status is the status of the response.
	Meta        GeneralResponse `json:"meta"`                   // Meta is the meta of the response.
	Data        any             `json:"data"`                   // Data is the data of the response.
	Total       int             `json:"total,omitempty"`        // Total is the total number of data items available. This field is optional and is typically used in paginated responses.
	PerPage     int             `json:"per_page,omitempty"`     // PerPage is the number of data items per page in a paginated response. This field is optional.
	CurrentPage int             `json:"current_page,omitempty"` // CurrentPage is the current page number in a paginated response. This field is optional.
	LastPage    int             `json:"last_page,omitempty"`    // LastPage is the number of the last page in a paginated response. This field is optional.`
	From        int             `json:"from,omitempty"`         // From is the starting index of the data items in the current page. This field is optional.
	To          int             `json:"to,omitempty"`           // To is the ending index of the data items in the current page. This field is optional.
}

// NewShipyardHandler initializes a new ShipyardHandler and sets up the routes for shipyard operations.
func NewShipyardHandler(e *echo.Echo, us shipyard.Usecase, middlewareAuth middleware.IMiddlewareAuth, mwP _reqContext.IMiddlewarePageRequest) {
	// Initialize a new ShipyardHandler with the provided use case and authentication middleware.
	handler = &ShipyardHandler{
		ShipyardUsecase: us,
		validator:       validator.New(),
		middlewareAuth:  middlewareAuth,
		mwPageRequest:   mwP,
	}

	// Create a new route group for shipyard operations.
	r := e.Group("v1/shipyard")
	// Use the authentication middleware for this route group.
	r.Use(handler.middlewareAuth.AuthorizationCheck)

	// Set up the routes for the shipyard operations.
	r.GET("/get-all", handler.FetchAllShipYards) // Fetch all shipyards.
	r.GET("/index", handler.FetchIndexShipYards, handler.mwPageRequest.PageRequestCtx)
	r.POST("", handler.StoreShipyard) // Store a new shipyard.

	// Create a new route group for operations that use a shipyard ID.
	rUseID := r.Group("/:id")
	// Use the contextShipyard middleware for this route group.
	rUseID.Use(contextShipyard)
	// Set up the routes for the shipyard operations that use a shipyard ID.
	rUseID.GET("", handler.ShowShipyardByIDOrCode) // Show a shipyard by its ID or code.
	rUseID.PUT("", handler.UpdateShipyard)         // Update a shipyard.
	rUseID.DELETE("", handler.DeleteShipyard)      // Delete a shipyard.
}

// returnOK is a function that returns a success response with the provided message and data.
// It returns an HTTP status code of 200 (OK) and a ReturnDataResponse object.
func returnOK(message string, data any, pdr *paginateDetailResult) (code int, i interface{}) {
	code = http.StatusOK
	returnData := ReturnDataResponse{
		Status: "success",
		Meta:   GeneralResponse{Message: message},
		Data:   data,
	}
	if pdr != nil {
		returnData.Total, returnData.PerPage, returnData.CurrentPage, returnData.LastPage, returnData.From, returnData.To = pdr.getPaginateDetail()
	}
	i = returnData
	return
}

// paginateDetailResult is a struct that represents the details of pagination.
// It contains total number of items, items per page, current page number, last page number, and the range (from, to) of items on the current page.
type paginateDetailResult struct {
	total       int // total is the total number of items.
	perPage     int // perPage is the number of items per page.
	lastPage    int // lastPage is the number of the last page.
	currentPage int // currentPage is the current page number.
	from        int // from is the starting index of the items on the current page.
	to          int // to is the ending index of the items on the current page.
}

// getPaginateDetail is a method for the paginateDetailResult struct.
// It returns the details of the pagination: total number of items, items per page, current page number, last page number, and the range (from, to) of items on the current page.
func (pdr *paginateDetailResult) getPaginateDetail() (
	total, perpage, currentpage, lastpage, from, to int) {
	total, perpage, currentpage, lastpage, from, to = pdr.total, pdr.perPage, pdr.currentPage, pdr.lastPage, pdr.from, pdr.to
	return
}

// returnBadRequest is a function that returns a bad request response with the provided message.
// It returns an HTTP status code of 400 (Bad Request) and a GeneralResponse object.
func returnBadRequest(message string) (code int, errResponse GeneralResponse) {
	code = http.StatusBadRequest
	errResponse = GeneralResponse{
		Message: message,
	}
	return
}

// FetchAllShipYards is a method of the ShipyardHandler struct.
// It fetches all active shipyards.
func (h *ShipyardHandler) FetchAllShipYards(c echo.Context) (err error) {
	// Call the FetchAllActiveShipyards method of the shipyard use case to fetch all active shipyards.
	data, _, _, err := h.ShipyardUsecase.FetchAllActiveShipyards()

	// If an error occurred, log it and return a bad request response.
	if err != nil {
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Return a success response with the fetched shipyards.
	err = c.JSON(returnOK("Success to retrieve shipyards", data, nil))
	return
}

func (h *ShipyardHandler) FetchIndexShipYards(c echo.Context) (err error) {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	// Call the FetchAllActiveShipyards method of the shipyard use case to fetch all active shipyards.
	data, total, lastpage, err := h.ShipyardUsecase.FetchShipyards(*pageRequest)

	// If an error occurred, log it and return a bad request response.
	if err != nil {
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Return a success response with the fetched shipyards.
	err = c.JSON(returnOK("Success to retrieve shipyards", data, &paginateDetailResult{
		total:       total,
		perPage:     pageRequest.PerPage,
		lastPage:    lastpage,
		currentPage: pageRequest.Page,
		from:        1,
		to:          lastpage,
	}))
	return
}

// ShowShipyardByIDOrCode is a method of the ShipyardHandler struct.
// It retrieves a shipyard by its ID or code.
func (h *ShipyardHandler) ShowShipyardByIDOrCode(c echo.Context) (err error) {
	// Get the shipyard data from the context.
	data := c.Get("shipyard").(models.Shipyard)

	// Return a success response with the retrieved shipyard.
	err = c.JSON(returnOK("Success to retrieve shipyard", data, nil))
	return
}

// StoreShipyard is a method of the ShipyardHandler struct.
// It stores a new shipyard.
func (h *ShipyardHandler) StoreShipyard(c echo.Context) (err error) {
	// Bind the request data to a ReqShipyard object.
	req := dto.ReqShipyard{}
	if err := c.Bind(&req); err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Validate the ReqShipyard object.
	if err := h.validator.Struct(req); err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Call the StoreShipyard method of the shipyard use case to store the new shipyard.
	result, err := h.ShipyardUsecase.StoreShipyard(req)
	if err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Return a success response with the stored shipyard.
	err = c.JSON(returnOK("Success to store shipyard", result, nil))
	return
}

// UpdateShipyard is a method of the ShipyardHandler struct.
// It updates an existing shipyard.
func (h *ShipyardHandler) UpdateShipyard(c echo.Context) (err error) {
	// Get the shipyard data from the context.
	data := c.Get("shipyard").(models.Shipyard)
	// Bind the request data to a ReqShipyard object.
	req := dto.ReqShipyard{}
	if err := c.Bind(&req); err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Validate the ReqShipyard object.
	if err := h.validator.Struct(req); err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// Call the UpdateShipyard method of the shipyard use case to update the shipyard.
	result, err := h.ShipyardUsecase.UpdateShipyard(data.ID, req)
	if err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Return a success response with the updated shipyard.
	err = c.JSON(returnOK("Success to update shipyard", result, nil))
	return
}

// DeleteShipyard is a method of the ShipyardHandler struct.
// It deletes a shipyard.
func (h *ShipyardHandler) DeleteShipyard(c echo.Context) (err error) {
	// Get the shipyard data from the context.
	data := c.Get("shipyard").(models.Shipyard)
	// Call the DeleteShipyard method of the shipyard use case to delete the shipyard.
	err = h.ShipyardUsecase.DeleteShipyard(data)
	if err != nil {
		// If an error occurred, log it and return a bad request response.
		zap.S().Error(err)
		return c.JSON(returnBadRequest(err.Error()))
	}

	// Return a success response.
	err = c.JSON(returnOK("Success to delete shipyard", nil, nil))
	return
}
