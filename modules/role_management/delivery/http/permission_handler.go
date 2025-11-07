package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

// permission scope
// get permission
// get index permission
// get all permission
// export all account based on type

func (handler *RoleManagementHandler) GetIndexPermission(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.RoleUseCase.GetIndexPermission(c, *pageRequest)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respPermission := []dto.RespPermissionIndex{}

	for _, v := range res {
		respPermission = append(respPermission, dto.ToRespPermissionIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respPermission, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *RoleManagementHandler) GetAllPermission(c echo.Context) error {

	res, err := handler.RoleUseCase.GetAllPermission(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respPermission := []dto.RespPermission{}

	for _, v := range res {
		respPermission = append(respPermission, dto.ToRespPermission(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respPermission)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetPermissionByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	res, err := handler.RoleUseCase.GetPermissionByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespPermissionDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetDuplicatedPermission(c echo.Context) error {
	req := new(dto.ReqCheckDuplicatedPermission)
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

	// PermissionId can be null
	if req.ExcludedPermissionId != uuid.Nil {
		uid = req.ExcludedPermissionId
	}

	res, err := handler.RoleUseCase.PermissionNameIsNotDuplicated(c, req.Name, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusNotFound, response.SetErrorResponse(http.StatusNotFound, "Permission Info with such name is not found"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Permission object
	resResp := dto.ToRespPermission(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
