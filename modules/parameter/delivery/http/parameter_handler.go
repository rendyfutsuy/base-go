package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/parameter"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type ParameterHandler struct {
	Usecase              parameter.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewParameterHandler(e *echo.Echo, uc parameter.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &ParameterHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	r := e.Group("/v1/parameter")
	r.Use(h.middlewareAuth.AuthorizationCheck)

	// Permissions
	// View:   parameter.view
	// Create: parameter.create
	// Update: parameter.update
	// Delete: parameter.delete
	// Export: parameter.export

	// Index with pagination + search
	r.GET("", h.GetIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation([]string{"parameter.view"}))

	// Export (no pagination, same filters) - must be before /:id to avoid route conflict
	r.GET("/export", h.Export, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation([]string{"parameter.export"}))

	// Get by ID (detail) - must be after /export to avoid route conflict
	r.GET("/:id", h.GetByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation([]string{"parameter.view"}))

	// Create
	r.POST("", h.Create, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation([]string{"parameter.create"}))

	// Update
	r.PUT("/:id", h.Update, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation([]string{"parameter.update"}))

	// Delete
	r.DELETE("/:id", h.Delete, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation([]string{"parameter.delete"}))
}

// Create godoc
// @Summary		Create a new parameter
// @Description	Create a new parameter with provided code, name, value, and description
// @Tags			Parameter
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateParameter	true	"Parameter creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespParameter}	"Successfully created parameter"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/parameter [post]
func (h *ParameterHandler) Create(c echo.Context) error {
	req := new(dto.ReqCreateParameter)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	user := c.Get("user")
	_ = user // not used; keep signature parity
	res, err := h.Usecase.Create(c, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespParameter(*res))
	return c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary		Update parameter
// @Description	Update an existing parameter's information
// @Tags			Parameter
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string					true	"Parameter UUID"
// @Param			request	body	dto.ReqUpdateParameter	true	"Updated parameter data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespParameter}	"Successfully updated parameter"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Parameter not found"
// @Router			/v1/parameter/{id} [put]
func (h *ParameterHandler) Update(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateParameter)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	res, err := h.Usecase.Update(c, id, req, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespParameter(*res))
	return c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary		Delete parameter
// @Description	Delete an existing parameter by ID
// @Tags			Parameter
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"Parameter UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted parameter"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Parameter not found"
// @Router			/v1/parameter/{id} [delete]
func (h *ParameterHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.Delete(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: "Successfully delete Parameter"})
	return c.JSON(http.StatusOK, resp)
}

// GetIndex godoc
// @Summary		Get list of parameters with pagination
// @Description	Retrieve a paginated list of parameters with optional filters
// @Tags			Parameter
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int						false	"Page number"
// @Param			per_page	query		int						false	"Items per page"
// @Param			search		query		string					false	"Search keyword"
// @Param			types		query		[]string				false	"Filter by parameter types (array)"
// @Param			names		query		[]string				false	"Filter by parameter names (array)"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespParameterIndex}	"Successfully retrieved parameters"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/parameter [get]
func (h *ParameterHandler) GetIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqParameterIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, total, err := h.Usecase.GetIndex(c, *pageRequest, *filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respParameter := []dto.RespParameterIndex{}

	for _, v := range res {
		respParameter = append(respParameter, dto.ToRespParameterIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respParameter, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

// GetByID godoc
// @Summary		Get parameter by ID
// @Description	Retrieve a single parameter by its ID
// @Tags			Parameter
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"Parameter UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespParameter}	"Successfully retrieved parameter"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Parameter not found"
// @Router			/v1/parameter/{id} [get]
func (h *ParameterHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespParameter(*res))
	return c.JSON(http.StatusOK, resp)
}

// Export godoc
// @Summary		Export parameters to Excel
// @Description	Export parameters to Excel file (.xlsx) with optional search and filter. Same search and filter logic as index but without pagination.
// @Tags			Parameter
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search		query		string					false	"Search keyword"
// @Param			types		query		[]string				false	"Filter by parameter types (array)"
// @Param			names		query		[]string				false	"Filter by parameter names (array)"
// @Success		200			{file}		binary	"Excel file with parameters data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/parameter/export [get]
func (h *ParameterHandler) Export(c echo.Context) error {
	// validate filter req.
	// initialize filter
	filter := new(dto.ReqParameterIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	excelBytes, err := h.Usecase.Export(c, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=parameters.xlsx")
	return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}
