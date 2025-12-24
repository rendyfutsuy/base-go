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
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/group"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

type GroupHandler struct {
	Usecase              group.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewGroupHandler(e *echo.Echo, uc group.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &GroupHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	r := e.Group("v1/group")
	r.Use(h.middlewareAuth.AuthorizationCheck)

	// Permissions
	// View:   group.view
	// Create: group.create
	// Update: group.update
	// Delete: group.delete
	// Export: group.export
	permissionToView := []string{"group.view"}
	permissionToCreate := []string{"group.create"}
	permissionToUpdate := []string{"group.update"}
	permissionToDelete := []string{"group.delete"}
	permissionToExport := []string{"group.export"}

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
// @Summary		Create a new group
// @Description	Create a new goods group with provided name
// @Tags			Golongan
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateGroup	true	"Group creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespGroup}	"Successfully created group"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/group [post]
func (h *GroupHandler) Create(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	req := new(dto.ReqCreateGroup)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	res, err := h.Usecase.Create(ctx, req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	// Newly created group is always deletable (not used in any sub-group yet)
	// Set Deletable to true since it's not included in Create query
	res.Deletable = true
	resp, _ = resp.SetResponse(dto.ToRespGroup(*res))
	return c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary		Update group
// @Description	Update an existing goods group's information
// @Tags			Golongan
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"Group UUID"
// @Param			request	body	dto.ReqUpdateGroup	true	"Updated group data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespGroup}	"Successfully updated group"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Group not found"
// @Router			/v1/group/{id} [put]
func (h *GroupHandler) Update(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	id := c.Param("id")
	req := new(dto.ReqUpdateGroup)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	res, err := h.Usecase.Update(ctx, id, req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	// deletable is already included in query
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespGroup(*res))
	return c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary		Delete group
// @Description	Delete an existing goods group by ID
// @Tags			Golongan
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string	true	"Group UUID"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully deleted group"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Group not found"
// @Router			/v1/group/{id} [delete]
func (h *GroupHandler) Delete(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// Get user ID from context
	user := c.Get("user")
	userID := ""
	if user != nil {
		if userModel, ok := user.(models.User); ok {
			userID = userModel.ID.String()
		}
	}

	id := c.Param("id")
	if err := h.Usecase.Delete(ctx, id, userID); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(Response{Message: constants.GroupDeleteSuccess})
	return c.JSON(http.StatusOK, resp)
}

// GetIndex godoc
// @Summary		Get list of groups with pagination
// @Description	Retrieve a paginated list of goods groups with optional filters
// @Tags			Golongan
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int					false	"Page number"
// @Param			per_page	query		int					false	"Items per page"
// @Param			search		query		string				false	"Search keyword"
// @Param			filter		query		dto.ReqGroupIndexFilter	false	"Filter options"
// @Success		200			{object}	response.PaginationResponse{data=[]dto.RespGroupIndex}	"Successfully retrieved groups"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/group [get]
func (h *GroupHandler) GetIndex(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	pageRequest := c.Get("page_request").(*request.PageRequest)

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqGroupIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, total, err := h.Usecase.GetIndex(ctx, *pageRequest, *filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respGroup := []dto.RespGroupIndex{}

	// Map groups to response with deletable status (deletable is already included in query)
	for _, v := range res {
		respGroup = append(respGroup, dto.ToRespGroupIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respGroup, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

// GetByID godoc
// @Summary		Get group by ID
// @Description	Retrieve a single goods group by its ID
// @Tags			Golongan
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path		string	true	"Group UUID"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespGroup}	"Successfully retrieved group"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Group not found"
// @Router			/v1/group/{id} [get]
func (h *GroupHandler) GetByID(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	id := c.Param("id")
	res, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	// deletable is already included in query
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespGroup(*res))
	return c.JSON(http.StatusOK, resp)
}

// Export godoc
// @Summary		Export groups to Excel
// @Description	Export goods groups to Excel file (.xlsx) with optional search and filter. Same search and filter logic as index but without pagination.
// @Tags			Golongan
// @Accept			json
// @Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security		BearerAuth
// @Param			filter		query		dto.ReqGroupIndexFilter	false	"Filter options"
// @Success		200			{file}		binary	"Excel file with groups data"
// @Failure		400			{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401			{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/group/export [get]
func (h *GroupHandler) Export(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqGroupIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	excelBytes, err := h.Usecase.Export(ctx, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	c.Response().Header().Set(echo.HeaderContentType, constants.ExcelContent)
	c.Response().Header().Set(constants.FieldContentDisposition, constants.ExcelContentDisposition("groups.xlsx"))
	return c.Blob(http.StatusOK, constants.ExcelContent, excelBytes)
}
