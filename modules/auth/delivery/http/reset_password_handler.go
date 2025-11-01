package http

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// ResetPasswordRequest godoc
// @Summary		Request a password reset email
// @Description	Sends a password reset email to the provided email address
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ReqResetPasswordRequest		true	"Reset Password Request"
// @Success		200		{object}	GeneralResponse{message=string}	"Successfully Send Reset Email Request"
// @Failure		400		{object}	GeneralResponse{message=string}	"Bad Request"
// @Router			/v1/auth/reset-password/request [post]
func (handler *AuthHandler) ResetPasswordRequest(c echo.Context) error {
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqResetPasswordRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// call update profile function
	err := handler.AuthUseCase.RequestResetPassword(ctx, req.Email)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Send Reset Email Request"})
}

// ResetUserPassword godoc
// @Summary		Reset user password
// @Description	Resets the user password using a valid token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			token	path		string							true	"Password Reset Token"
// @Param			request	body		dto.ReqResetPassword			true	"Reset User Password"
// @Success		200		{object}	GeneralResponse{message=string}	"Successfully Reset Password"
// @Failure		400		{object}	GeneralResponse{message=string}	"Bad Request"
// @Router			/v1/auth/reset-password/request/{token} [post]
func (handler *AuthHandler) ResetUserPassword(c echo.Context) error {
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqResetPassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// get token from route
	token := c.Param("token")

	// Decode the Base64 token
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: "Invalid token"})
	}

	//get user through password reset token
	err = handler.AuthUseCase.ResetUserPassword(ctx, req.Password, string(decodedToken))

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Reset Password"})
}
