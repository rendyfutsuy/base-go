package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/rendyfutsuy/base-go/helpers/middleware"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// GeneralResponse represent the response error struct
type GeneralResponse struct {
	Message string `json:"message"`
}

type ResponseAuth struct {
	AccessToken string `json:"access_token"`
}

type ResponseError struct {
	Message string `json:"message"`
}

// AuthHandler represent the http handler for auth
type AuthHandler struct {
	AuthUseCase    auth.Usecase
	validator      *validator.Validate
	middlewareAuth middleware.IMiddlewareAuth
}

// NewAuthHandler will initialize the auth/ resources endpoint
func NewAuthHandler(e *echo.Echo, us auth.Usecase, middlewareAuth middleware.IMiddlewareAuth) {
	handler := &AuthHandler{
		AuthUseCase:    us,
		validator:      validator.New(),
		middlewareAuth: middlewareAuth,
	}

	r := e.Group("v1/auth")

	// not using middleware
	r.POST("/login",
		handler.Authenticate,
	)

	r.POST("/reset-password/request",
		handler.ResetPasswordRequest,
	)

	r.POST("/reset-password/request/:token",
		handler.ResetUserPassword,
	)

	// use middleware
	r.Use(handler.middlewareAuth.AuthorizationCheck)
	r.POST("/logout",
		handler.SignOut,
		handler.middlewareAuth.AuthorizationCheck,
	)

	r.GET("/profile",
		handler.GetProfile,
		handler.middlewareAuth.AuthorizationCheck,
	)

	r.POST("/profile",
		handler.UpdateProfile,
		handler.middlewareAuth.AuthorizationCheck,
	)

	r.POST("/profile/my-password",
		handler.UpdateMyPassword,
		handler.middlewareAuth.AuthorizationCheck,
	)
}

// @Summary		Authenticate user
// @Description	Authenticates a user and returns an access token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ReqAuthUser	true	"User login and password"
// @Success		200		{object}	ResponseAuth	"Successfully authenticated"
// @Failure		400		{object}	GeneralResponse	"Invalid request"
// @Failure		419		{object}	GeneralResponse	"User password expired"
// @Router			/v1/auth/login [post]
func (handler *AuthHandler) Authenticate(c echo.Context) error {
	// Validate input
	req := new(dto.ReqAuthUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// Assuming the Authenticate method on AuthUseCase does the actual authentication
	token, err := handler.AuthUseCase.Authenticate(c, req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseAuth{AccessToken: token})
}

// @Summary		Sign out user
// @Description	Logs out the user by invalidating the session token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	GeneralResponse	"Successfully logged out"
// @Failure		400	{object}	GeneralResponse	"Logout failed"
// @Failure		401	{object}	GeneralResponse	"Unauthorized"
// @Router			/v1/auth/logout [post]
func (handler *AuthHandler) SignOut(c echo.Context) error {

	// parse token
	token := c.Get("token").(string)

	// initiate session destroy
	err := handler.AuthUseCase.SignOut(c, token)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Logged Out"})
}

// @Summary		Get user profile
// @Description	Retrieves the profile of the authenticated user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	dto.UserProfile	"User profile data"
// @Failure		400	{object}	GeneralResponse	"Failed to get profile"
// @Failure		401	{object}	GeneralResponse	"Unauthorized"
// @Router			/v1/auth/profile [get]
func (handler *AuthHandler) GetProfile(c echo.Context) error {

	// initiate session destroy
	user, err := handler.AuthUseCase.GetProfile(c)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary		Update user profile
// @Description	Updates the profile information of the authenticated user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqUpdateProfile	true	"Updated user profile data"
// @Success		200		{object}	GeneralResponse			"Successfully updated profile"
// @Failure		400		{object}	GeneralResponse			"Invalid request"
// @Failure		401		{object}	GeneralResponse			"Unauthorized"
// @Router			/v1/auth/profile [put]
func (handler *AuthHandler) UpdateProfile(c echo.Context) error {

	// Validate input
	req := new(dto.ReqUpdateProfile)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// call update profile function
	err := handler.AuthUseCase.UpdateProfile(c, *req)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Updated Profile"})
}

// @Summary		Update user password
// @Description	Updates the password of the authenticated user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqUpdatePassword	true	"New password data"
// @Success		200		{object}	GeneralResponse			"Successfully updated password"
// @Failure		400		{object}	GeneralResponse			"Invalid request"
// @Failure		401		{object}	GeneralResponse			"Unauthorized"
// @Router			/v1/auth/password [put]
func (handler *AuthHandler) UpdateMyPassword(c echo.Context) error {
	// Validate input
	req := new(dto.ReqUpdatePassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// call update profile function
	err := handler.AuthUseCase.UpdateMyPassword(c, *req)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: "Successfully Updated My Password"})
}
