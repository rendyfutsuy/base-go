package http

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/auth/dto"
)

func (handler *AuthHandler) ResetPasswordRequest(c echo.Context) error {
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
	err := handler.AuthUseCase.RequestResetPassword(c, req.Email)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Send Reset Email Request"})
}

func (handler *AuthHandler) ResetUserPassword(c echo.Context) error {
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
	err = handler.AuthUseCase.ResetUserPassword(c, req.Password, string(decodedToken))

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Reset Password"})
}
