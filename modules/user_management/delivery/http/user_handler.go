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

// CreateUser godoc
// @Summary		Create a new user
// @Description	Create a new user with provided information
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqCreateUser	true	"User creation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUser}	"Successfully created user"
// @Failure		400		{object}	ResponseError	"Bad request - validation error"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Router			/v1/user-management/user [post]
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

// GetIndexUser godoc
// @Summary		Get paginated list of users
// @Description	Retrieve a paginated list of users with optional filtering
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			page		query		int		false	"Page number"	default(1)
// @Param			per_page	query		int		false	"Items per page"	default(10)
// @Param			sort_by		query		string	false	"Sort column"
// @Param			sort_order	query		string	false	"Sort order (asc/desc)"
// @Param			search		query		string	false	"Search query"
// @Param			filter		query		dto.ReqUserIndexFilter	false	"Filter options"
// @Success		200		{object}	response.PaginationResponse{data=[]dto.RespUserIndex}	"Successfully retrieved users"
// @Failure		400		{object}	ResponseError	"Bad request"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Router			/v1/user-management/user [get]
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

	res, total, err := handler.UserUseCase.GetIndexUser(c, *pageRequest, *filter)

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

// GetAllUser godoc
// @Summary		Get all users
// @Description	Retrieve all users without pagination
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200		{object}	response.NonPaginationResponse{data=[]dto.RespUser}	"Successfully retrieved all users"
// @Failure		400		{object}	ResponseError	"Bad request"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Router			/v1/user-management/user/all [get]
func (handler *UserManagementHandler) GetAllUser(c echo.Context) error {

	res, err := handler.UserUseCase.GetAllUser(c)
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

// GetUserByID godoc
// @Summary		Get user by ID
// @Description	Retrieve a specific user by their ID
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path		string	true	"User UUID"
// @Success		200	{object}	response.NonPaginationResponse{data=dto.RespUserDetail}	"Successfully retrieved user"
// @Failure		400	{object}	ResponseError	"Bad request - invalid UUID"
// @Failure		401	{object}	ResponseError	"Unauthorized"
// @Failure		404	{object}	ResponseError	"User not found"
// @Router			/v1/user-management/user/{id} [get]
func (handler *UserManagementHandler) GetUserByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	res, err := handler.UserUseCase.GetUserByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Summary		Update user
// @Description	Update an existing user's information
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"User UUID"
// @Param			request	body	dto.ReqUpdateUser	true	"Updated user data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUser}	"Successfully updated user"
// @Failure		400		{object}	ResponseError	"Bad request - validation error"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Failure		404		{object}	ResponseError	"User not found"
// @Router			/v1/user-management/user/{id} [put]
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

	res, err := handler.UserUseCase.UpdateUser(c, id, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// GetDuplicatedUser godoc
// @Summary		Check if user name is duplicated
// @Description	Check if a user name already exists in the database
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body	dto.ReqCheckDuplicatedUser	true	"Check duplicated user request"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUser}	"User with such name exists"
// @Failure		400		{object}	ResponseError	"Bad request"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Failure		404		{object}	ResponseError	"User with such name is not found"
// @Router			/v1/user-management/user/check-name [get]
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

	res, err := handler.UserUseCase.UserNameIsNotDuplicated(c, req.FullName, uid)

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

// BlockUser godoc
// @Summary		Block a user
// @Description	Block a user account and revoke all their tokens
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"User UUID"
// @Param			request	body	dto.ReqBlockUser	true	"Block user request"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUserDetail}	"Successfully blocked user"
// @Failure		400		{object}	ResponseError	"Bad request"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Failure		404		{object}	ResponseError	"User not found"
// @Router			/v1/user-management/user/{id}/block [patch]
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

// ActivateUser godoc
// @Summary		Activate a user
// @Description	Activate or update status of a user account
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string				true	"User UUID"
// @Param			request	body	dto.ReqActivateUser	true	"Activate user request"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUserDetail}	"Successfully activated user"
// @Failure		400		{object}	ResponseError	"Bad request"
// @Failure		401		{object}	ResponseError	"Unauthorized"
// @Failure		404		{object}	ResponseError	"User not found"
// @Router			/v1/user-management/user/{id}/assign-status [patch]
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
