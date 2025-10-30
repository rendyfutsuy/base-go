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
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
)

// user scope
// create user
// update user
// delete user
// get user
// get index user
// get all user

func (handler *UserManagementHandler) CreateUser(c echo.Context) error {

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	fmt.Println(authId)

	req := new(dto.ReqCreateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.UserUseCase.CreateUser(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) GetIndexUser(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqUserIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, total, err := handler.UserUseCase.GetIndexUser(*pageRequest, *filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	respUser := []dto.RespUserIndex{}

	for _, v := range res {
		respUser = append(respUser, dto.ToRespUserIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respUser, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *UserManagementHandler) GetAllUser(c echo.Context) error {

	res, err := handler.UserUseCase.GetAllUser()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	respUser := []dto.RespUser{}

	for _, v := range res {
		respUser = append(respUser, dto.ToRespUser(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respUser)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) GetUserByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	res, err := handler.UserUseCase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) UpdateUser(c echo.Context) error {

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	req := new(dto.ReqUpdateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.UserUseCase.UpdateUser(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) GetDuplicatedUser(c echo.Context) error {
	req := new(dto.ReqCheckDuplicatedUser)
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

	// UserId can be null
	if req.ExcludedUserId != uuid.Nil {
		uid = req.ExcludedUserId
	}

	res, err := handler.UserUseCase.UserNameIsNotDuplicated(req.FullName, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: "User Info with such name is not found"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// if name already uses by existing account info, return User object
	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) BlockUser(c echo.Context) error {
	req := new(dto.ReqBlockUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	// get Block User
	// add revoke all user auth token
	res, err := handler.UserUseCase.BlockUser(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// revoke user token

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *UserManagementHandler) ActivateUser(c echo.Context) error {
	req := new(dto.ReqActivateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	// get Active User
	// add revoke all user auth token
	res, err := handler.UserUseCase.ActivateUser(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// revoke user token

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
