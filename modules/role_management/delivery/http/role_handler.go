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

func (handler *RoleManagementHandler) CreateRole(c echo.Context) error {

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	fmt.Println(authId)

	req := new(dto.ReqCreateRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.RoleUseCase.CreateRole(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetIndexRole(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.RoleUseCase.GetIndexRole(c, *pageRequest)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	respRole := []dto.RespRoleIndex{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRoleIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respRole, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *RoleManagementHandler) GetAllRole(c echo.Context) error {

	res, err := handler.RoleUseCase.GetAllRole(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	respRole := []dto.RespRole{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRole(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respRole)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetRoleByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	res, err := handler.RoleUseCase.GetRoleByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRoleDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) UpdateRole(c echo.Context) error {

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	req := new(dto.ReqUpdateRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.RoleUseCase.UpdateRole(c, id, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) DeleteRole(c echo.Context) error {
	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	res, err := handler.RoleUseCase.SoftDeleteRole(c, id, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetDuplicatedRole(c echo.Context) error {
	req := new(dto.ReqCheckDuplicatedRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate input
	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// initialize uid
	uid := uuid.Nil

	// RoleId can be null
	if req.ExcludedRoleId != uuid.Nil {
		uid = req.ExcludedRoleId
	}

	res, err := handler.RoleUseCase.RoleNameIsNotDuplicated(c, req.Name, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: "Role Info with such name is not found"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// if name already uses by existing account info, return Role object
	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetMyPermissions(c echo.Context) error {
	token := c.Get("token").(string)

	res, err := handler.RoleUseCase.MyPermissionsByUserToken(c, token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRoleDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
