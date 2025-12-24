package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
)

// user_password scope
// Update user Password
// Verify User Password

// UpdateUserPassword godoc
// @Summary		Update user password
// @Description	Update a user's password by their ID
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			id		path	string					true	"User UUID"
// @Param			request	body	dto.ReqUpdateUserPassword	true	"Password update data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUser}	"Successfully updated password"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Failure		404		{object}	response.NonPaginationResponse	"User not found"
// @Router			/v1/user-management/user/{id}/password [patch]
func (handler *UserManagementHandler) UpdateUserPassword(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get params ID
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	req := new(dto.ReqUpdateUserPassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// update Update User Password By Passed ID
	err = handler.UserUseCase.UpdateUserPassword(ctx, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// get User ID By Passed ID
	res, err := handler.UserUseCase.GetUserByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// ConfirmCurrentUserPassword godoc
// @Summary		Confirm current user password
// @Description	Verify the current authenticated user's password
// @Tags			User Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body	dto.ReqConfirmationUserPassword	true	"Password confirmation data"
// @Success		200		{object}	response.NonPaginationResponse{data=dto.RespUser}	"Password confirmed successfully"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - invalid password"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/user-management/user/password-confirmation [post]
func (handler *UserManagementHandler) ConfirmCurrentUserPassword(c echo.Context) error {
	// initialize context from echo
	ctx := c.Request().Context()

	// get auth ID
	user := c.Get("user")
	authId := user.(models.User).ID.String()

	req := new(dto.ReqConfirmationUserPassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// update Update User Password By Passed ID
	err := handler.UserUseCase.AssertCurrentUserPassword(ctx, authId, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// get User ID By Passed ID
	res, err := handler.UserUseCase.GetUserByID(ctx, authId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespUser(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
