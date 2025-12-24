package http

import (
	"fmt"
	"net/http"
	"strings"

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
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/user-management/user [post]
func (handler *UserManagementHandler) CreateUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	fmt.Println(authId)

	req := new(dto.ReqCreateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, err := handler.UserUseCase.CreateUser(ctx, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
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
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/user-management/user [get]
func (handler *UserManagementHandler) GetIndexUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	pageRequest := c.Get("page_request").(*request.PageRequest)

	// validate filter req.
	// initialize filter
	filter := new(dto.ReqUserIndexFilter)

	// Bind form-data to the DTO
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Validate the request if necessary
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, total, err := handler.UserUseCase.GetIndexUser(ctx, *pageRequest, *filter)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respUser := []dto.RespUserIndex{}

	for _, v := range res {
		respUser = append(respUser, dto.ToRespUserIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respUser, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
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
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/user-management/user/all [get]
func (handler *UserManagementHandler) GetAllUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	res, err := handler.UserUseCase.GetAllUser(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
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
// @Failure		400	{object}	response.NonPaginationResponse	"Bad request - invalid UUID"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404	{object}	response.NonPaginationResponse	"User not found"
// @Router			/v1/user-management/user/{id} [get]
func (handler *UserManagementHandler) GetUserByID(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	res, err := handler.UserUseCase.GetUserByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
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
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"User not found"
// @Router			/v1/user-management/user/{id} [put]
func (handler *UserManagementHandler) UpdateUser(c echo.Context) error {
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

	req := new(dto.ReqUpdateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// update user data
	res, err := handler.UserUseCase.UpdateUser(ctx, id, req, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// update user password, only if you are super admin.
	// if not skip this process entirely.
	if req.Password != "" &&
		strings.EqualFold(user.(models.User).RoleName, constants.AuthRoleSuperAdmin) {
		// apped password to update password validation
		reqPassword := new(dto.ReqUpdateUserPassword)
		reqPassword.NewPassword = req.Password
		reqPassword.PasswordConfirmation = req.PasswordConfirmation

		// update Update User Password By Passed ID
		err = handler.UserUseCase.UpdateUserPasswordNoCheckRequired(ctx, id, reqPassword)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
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
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		200		{object}	response.NonPaginationResponse	"User with such name is not found"
// @Router			/v1/user-management/user/check-name [post]
func (handler *UserManagementHandler) GetDuplicatedUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	req := new(dto.ReqCheckDuplicatedUser)
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

	// UserId can be null
	if req.ExcludedUserId != uuid.Nil {
		uid = req.ExcludedUserId
	}

	res, err := handler.UserUseCase.UserNameIsNotDuplicated(ctx, req.UserName, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusOK, response.SetErrorResponse(http.StatusOK, "User Info with such name is not found"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return User object
	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)
	resp.Status = http.StatusConflict

	return c.JSON(http.StatusConflict, resp)
}

// 2025/11/04: unused - commented first
// // BlockUser godoc
// // @Summary		Block a user
// // @Description	Block a user account and revoke all their tokens
// // @Tags			User Management
// // @Accept			json
// // @Produce		json
// // @Security		BearerAuth
// // @Param			id		path	string				true	"User UUID"
// // @Param			request	body	dto.ReqBlockUser	true	"Block user request"
// // @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUserDetail}	"Successfully blocked user"
// // @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// // @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// // @Failure		404		{object}	response.NonPaginationResponse	"User not found"
// // @Router			/v1/user-management/user/{id}/block [patch]
func (handler *UserManagementHandler) BlockUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	req := new(dto.ReqBlockUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	// get Block User
	// add revoke all user auth token
	res, err := handler.UserUseCase.BlockUser(ctx, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// revoke user token

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// 2025/11/04: unused - commented first
// // ActivateUser godoc
// // @Summary		Activate a user
// // @Description	Activate or update status of a user account
// // @Tags			User Management
// // @Accept			json
// // @Produce		json
// // @Security		BearerAuth
// // @Param			id		path	string				true	"User UUID"
// // @Param			request	body	dto.ReqActivateUser	true	"Activate user request"
// // @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUserDetail}	"Successfully activated user"
// // @Failure		400		{object}	response.NonPaginationResponse	"Bad request"
// // @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// // @Failure		404		{object}	response.NonPaginationResponse	"User not found"
// // @Router			/v1/user-management/user/{id}/assign-status [patch]
func (handler *UserManagementHandler) ActivateUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	req := new(dto.ReqActivateUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	// get Active User
	// add revoke all user auth token
	res, err := handler.UserUseCase.ActivateUser(ctx, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// revoke user token

	resResp := dto.ToRespUserDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary		Soft delete a user
// @Description	Soft delete a user account by setting deleted_at timestamp. The user will be marked as deleted but remain in the database. Requires 'api.user-management.user.delete' permission.
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id	path	string	true	"User UUID"
// @Success		200	{object}	response.NonPaginationResponse{data=dto.RespUser}	"Successfully soft deleted user"
// @Failure		400	{object}	response.NonPaginationResponse	"Bad request - invalid UUID or user not found"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		403	{object}	response.NonPaginationResponse	"Forbidden - insufficient permissions"
// @Failure		404	{object}	response.NonPaginationResponse	"User not found"
// @Router			/v1/user-management/user/{id} [delete]
func (handler *UserManagementHandler) DeleteUser(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	res, err := handler.UserUseCase.SoftDeleteUser(ctx, id, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
