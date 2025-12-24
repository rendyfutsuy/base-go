package http

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

// role scope
// create role
// update role
// delete role
// get role
// get index role
// get all role

// CreateRole godoc
// @Summary		Create a new role
// @Description	Create a new role with provided information
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateRole	true	"Role creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespRole}	"Successfully created role"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/role-management/role [post]
func (handler *RoleManagementHandler) CreateRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	fmt.Println(authId)

	req := new(dto.ReqCreateRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, err := handler.RoleUseCase.CreateRole(ctx, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// GetIndexRole godoc
// @Summary		Get paginated list of roles
// @Description	Retrieve a paginated list of roles with optional filtering
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int		false	"Page number"	default(1)
// @Param			per_page	query		int		false	"Items per page"	default(10)
// @Param			sort_by		query		string	false	"Sort column"
// @Param			sort_order	query		string	false	"Sort order (asc/desc)"
// @Param			search		query		string	false	"Search query"
// @Success		200		{object}	response.PaginationResponse{data=[]dto.RespRoleIndex}	"Successfully retrieved roles"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/role-management/role [get]
func (handler *RoleManagementHandler) GetIndexRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.RoleUseCase.GetIndexRole(ctx, *pageRequest)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respRole := []dto.RespRoleIndex{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRoleIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respRole, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

// GetAllRole godoc
// @Summary		Get all roles
// @Description	Retrieve all roles without pagination
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200		{object}	response.NonPaginationResponse{data=[]dto.RespRole}	"Successfully retrieved all roles"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/role-management/role/all [get]
func (handler *RoleManagementHandler) GetAllRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	res, err := handler.RoleUseCase.GetAllRole(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respRole := []dto.RespRole{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRole(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respRole)

	return c.JSON(http.StatusOK, resp)
}

// GetRoleByID godoc
// @Summary		Get role by ID
// @Description	Retrieve a specific role by its ID
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"Role UUID"
// @Success		200	{object}	response.NonPaginationResponse{data=dto.RespRoleDetail}	"Successfully retrieved role"
// @Failure		400	{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404	{object}	response.NonPaginationResponse	"Role not found"
// @Router			/v1/role-management/role/{id} [get]
func (handler *RoleManagementHandler) GetRoleByID(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	res, err := handler.RoleUseCase.GetRoleByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Role object
	modules := []dto.RespPermissionGroupByModule{}

	modules, err = handler.buildPermissionGroupsByModule(ctx, id)
	if err != nil {
		// Handle specific error for UUID parsing
		if id != "" {
			if _, parseErr := uuid.Parse(id); parseErr != nil {
				return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
			}
		}
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespRoleDetail(*res, modules)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// UpdateRole godoc
// @Summary		Update role
// @Description	Update an existing role's information
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"Role UUID"
// @Param			request	body	dto.ReqUpdateRole	true	"Updated role data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespRole}	"Successfully updated role"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"Role not found"
// @Router			/v1/role-management/role/{id} [put]
func (handler *RoleManagementHandler) UpdateRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	req := new(dto.ReqUpdateRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, err := handler.RoleUseCase.UpdateRole(ctx, id, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// DeleteRole godoc
// @Summary		Delete role (soft delete)
// @Description	Soft delete a role by its ID
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path	string	true	"Role UUID"
// @Success		200	{object}	response.NonPaginationResponse{data=dto.RespRole}	"Successfully deleted role"
// @Failure		400	{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404	{object}	response.NonPaginationResponse	"Role not found"
// @Router			/v1/role-management/role/{id} [delete]
func (handler *RoleManagementHandler) DeleteRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	res, err := handler.RoleUseCase.SoftDeleteRole(ctx, id, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if total user is greater than 0, return error
	// role that have been uses can not be deleted
	if res.TotalUser > 0 {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorCannotDeleteRole))
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// GetDuplicatedRole godoc
// @Summary		Check if role name is duplicated
// @Description	Check if a role name already exists in the database
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body	dto.ReqCheckDuplicatedRole	true	"Check duplicated role request"
// @Success		409		{object}	response.NonPaginationResponse{data=dto.RespRole}	"Role with such name exists"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		200		{object}	response.NonPaginationResponse	"Role with such name is not found"
// @Router			/v1/role-management/role/check-name [post]
func (handler *RoleManagementHandler) GetDuplicatedRole(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	req := new(dto.ReqCheckDuplicatedRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate input
	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// initialize uid
	uid := uuid.Nil

	// RoleId can be null
	if req.ExcludedRoleId != uuid.Nil {
		uid = req.ExcludedRoleId
	}

	res, err := handler.RoleUseCase.RoleNameIsNotDuplicated(ctx, req.Name, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusOK, response.SetErrorResponse(http.StatusOK, "Role Info with such name is not found"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Role object
	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)
	resp.Status = http.StatusConflict

	return c.JSON(http.StatusConflict, resp)
}

func (handler *RoleManagementHandler) GetMyPermissions(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	token := c.Get("token").(string)
	res, err := handler.RoleUseCase.MyPermissionsByUserToken(ctx, token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Role object
	modules := []dto.RespPermissionGroupByModule{}
	modules, err = handler.buildPermissionGroupsByModule(ctx, res.ID.String())
	if err != nil {
		// Handle specific error for UUID parsing
		if res.ID.String() != "" {
			if _, parseErr := uuid.Parse(res.ID.String()); parseErr != nil {
				return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
			}
		}
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespRoleDetail(*res, modules)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
