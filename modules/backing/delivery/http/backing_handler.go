package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/backing"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type BackingHandler struct {
	Usecase              backing.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewBackingHandler(e *echo.Echo, uc backing.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &BackingHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	r := e.Group("v1/backing")
	r.Use(h.middlewareAuth.AuthorizationCheck)

	// Permissions
	// View:   backing.view
	// Create: backing.create
	// Update: backing.update
	// Delete: backing.delete
	// Export: backing.export
	permissionToView := []string{"backing.view"}
	permissionToCreate := []string{"backing.create"}
	permissionToUpdate := []string{"backing.update"}
	permissionToDelete := []string{"backing.delete"}
	permissionToExport := []string{"backing.export"}

	// Index with pagination + search
	r.GET("", h.GetIndex, middleware.RequireActivatedUser, h.mwPageRequest.PageRequestCtx, h.middlewarePermission.PermissionValidation(permissionToView))

	// Export (no pagination, same filters) - must be before /:id to avoid route conflict
	r.GET("/export", h.Export, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToExport))

	// Get by ID (detail) - must be after /export to avoid route conflict
	r.GET("/:id", h.GetByID, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToView))

	// Create
	r.POST("", h.Create, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))

	// Update
	r.PUT("/:id", h.Update, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))

	// Delete
	r.DELETE("/:id", h.Delete, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))
}

// Create godoc
// @Summary		Create a new backing
// @Description	Create a new backing with provided type_id and name
// @Tags			Backing
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateBacking	true	"Backing creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespBacking}	"Successfully created backing"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/backing [post]
func (h *BackingHandler) Create(c echo.Context) error {
	req := new(dto.ReqCreateBacking)
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
	resp, _ = resp.SetResponse(dto.ToRespBacking(*res))
	return c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary		Update backing
// @Description	Update an existing backing's information
// @Tags			Backing
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"Backing UUID"
// @Param			request	body	dto.ReqUpdateBacking	true	"Updated backing data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespBacking}	"Successfully updated backing"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Backing not found"
// @Router			/v1/backing/{id} [put]
func (h *BackingHandler) Update(c echo.Context) error {
	id := c.Param("id")
	req := new(dto.ReqUpdateBacking)
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
	resp, _ = resp.SetResponse(dto.ToRespBacking(*res))
	return c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary		Delete backing
// @Description	Delete an existing backing by ID
// @Tags			Backing
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"Backing UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted backing"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Backing not found"
// @Router			/v1/backing/{id} [delete]
func (h *BackingHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.Usecase.Delete(c, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.BackingDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetIndex godoc
// @Summary		Get list of backings with pagination
// @Description	Retrieve a paginated list of backings with optional search and filters. Supports multiple filter values for backing_code, name, and type_id.
// @Tags			Backing
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page			query		int						false	"Page number"
// @Param			per_page		query		int						false	"Items per page"
// @Param			sort_by			query		string					false	"Sort column (allowed: id, type_id, backing_code, name, created_at, updated_at)"
// @Param			sort_order		query		string					false	"Sort order (asc or desc)"
// @Param			search			query		string					false	"Search keyword (searches in backing_code and name)"
// @Param			backing_codes	query		[]string				false	"Filter by backing codes (multiple values)"
// @Param			names			query		[]string				false	"Filter by names (multiple values)"
// @Param			type_ids		query		[]string				false	"Filter by type IDs (multiple values, UUID format)"
// @Param			filter			query		dto.ReqBackingIndexFilter	false	"Filter options"
// @Success		200				{object}	response.PaginationResponse{data=[]dto.RespBackingIndex}	"Successfully retrieved backings"
// @Failure		400				{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401				{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/backing [get]
func (h *BackingHandler) GetIndex(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqBackingIndexFilter)

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

	respBacking := []dto.RespBackingIndex{}

	for _, v := range res {
		respBacking = append(respBacking, dto.ToRespBackingIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respBacking, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

// GetByID godoc
// @Summary		Get backing by ID
// @Description	Retrieve a single backing by its ID
// @Tags			Backing
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"Backing UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespBacking}	"Successfully retrieved backing"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Backing not found"
// @Router			/v1/backing/{id} [get]
func (h *BackingHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.Usecase.GetByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespBacking(*res))
	return c.JSON(http.StatusOK, resp)
}

// Export godoc
// @Summary		Export backings to Excel
// @Description	Export backings to Excel file (.xlsx) with optional search and filter. Same search and filter logic as index but without pagination. Supports multiple filter values for backing_code, name, and type_id.
// @Tags			Backing
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			search			query		string					false	"Search keyword (searches in backing_code and name)"
// @Param			backing_codes	query		[]string				false	"Filter by backing codes (multiple values)"
// @Param			names			query		[]string				false	"Filter by names (multiple values)"
// @Param			type_ids		query		[]string				false	"Filter by type IDs (multiple values, UUID format)"
// @Param			filter			query		dto.ReqBackingIndexFilter	false	"Filter options"
// @Success		200				{file}		binary	"Excel file with backings data"
// @Failure		400				{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401				{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/backing/export [get]
func (h *BackingHandler) Export(c echo.Context) error {
	// validate filter req.
	// initialize filter
	filter := new(dto.ReqBackingIndexFilter)

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

	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("backings.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}
