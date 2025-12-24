package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// GeneralResponse represent the response error struct
type GeneralResponse struct {
	Message string `json:"message"`
}

type ResponseAuth struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IsFirstTimeLogin bool   `json:"is_first_time_login"`
}

type ResponseError struct {
	Message string `json:"message"`
}

// AuthHandler represent the http handler for auth
type AuthHandler struct {
	AuthUseCase    auth.Usecase
	validator      *validator.Validate
	middlewareAuth middleware.IMiddlewareAuth
	mwPageRequest  _reqContext.IMiddlewarePageRequest
}

// NewAuthHandler will initialize the auth/ resources endpoint
func NewAuthHandler(e *echo.Echo, us auth.Usecase, middlewareAuth middleware.IMiddlewareAuth, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &AuthHandler{
		AuthUseCase:    us,
		validator:      validator.New(),
		middlewareAuth: middlewareAuth,
		mwPageRequest:  mwP,
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

	r.POST("/refresh-token",
		handler.RefreshToken,
	)
}

// @Summary		Authenticate user
// @Description	Authenticates a user and returns an access token and is_first_time_login status
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ReqAuthUser	true	"User login and password"
// @Success		200		{object}	response.NonPaginationResponse{data=ResponseAuth}	"Successfully authenticated"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized - invalid credentials or user not found"
// @Router			/v1/auth/login [post]
func (handler *AuthHandler) Authenticate(c echo.Context) error {
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqAuthUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, constants.AuthUsernamePasswordNotFound))
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, constants.AuthUsernamePasswordNotFound))
	}

	// Authenticate user
	result, err := handler.AuthUseCase.Authenticate(ctx, req.Login, req.Password)
	if err != nil {
		// All other errors return 401
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(ResponseAuth{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		IsFirstTimeLogin: result.IsFirstTimeLogin,
	})

	return c.JSON(http.StatusOK, resp)
}

// @Summary		Sign out user
// @Description	Logs out the user by invalidating the session token
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	response.NonPaginationResponse	"Successfully logged out"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/auth/logout [post]
func (handler *AuthHandler) SignOut(c echo.Context) error {
	ctx := c.Request().Context()

	// parse token
	token := c.Get("token").(string)

	// initiate session destroy
	err := handler.AuthUseCase.SignOut(ctx, token)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(nil)
	resp.Message = constants.AuthLogoutSuccess
	return c.JSON(http.StatusOK, resp)
}

// @Summary		Get user profile
// @Description	Retrieves the profile of the authenticated user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	response.NonPaginationResponse{data=dto.UserProfile}	"User profile data"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/auth/profile [get]
func (handler *AuthHandler) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// parse token
	token := c.Get("token").(string)

	// initiate session destroy
	user, err := handler.AuthUseCase.GetProfile(ctx, token)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// Convert models.User to dto.UserProfile
	profile := dto.UserProfile{
		UserId:           user.ID.String(),
		Name:             user.FullName,
		Username:         user.Username,
		Email:            user.Email,
		IsFirstTimeLogin: user.IsFirstTimeLogin,
		Role:             user.RoleName,
		Permissions:      user.Permissions,
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(profile)
	return c.JSON(http.StatusOK, resp)
}

// 2025/11/04: unused - commented first
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
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqUpdateProfile)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// parse user ID from context
	userId := c.Get("userId").(string)

	// call update profile function
	err := handler.AuthUseCase.UpdateProfile(ctx, *req, userId)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(GeneralResponse{Message: constants.AuthProfileUpdated})
	return c.JSON(http.StatusOK, resp)
}

// @Summary		Update user password
// @Description	Updates the password of the authenticated user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			request	body		dto.ReqUpdatePassword	true	"New password data"
// @Success		200		{object}	response.NonPaginationResponse	"Successfully updated password"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/auth/profile/my-password [put]
func (handler *AuthHandler) UpdateMyPassword(c echo.Context) error {
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqUpdatePassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// parse user ID from context
	userId := c.Get("userId").(string)

	// call update profile function
	err := handler.AuthUseCase.UpdateMyPassword(ctx, *req, userId)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(nil)
	resp.Message = constants.AuthPasswordUpdated
	return c.JSON(http.StatusOK, resp)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Refresh access token
// @Description	Generates a new access token based on the provided bearer token. Access tokens that have expired can still be refreshed as long as the Redis session (TTL) has not expired. If the token is revoked or the Redis session has expired, returns an error.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	response.NonPaginationResponse{data=ResponseAuth}	"Successfully refreshed token"
// @Failure		401	{object}	response.NonPaginationResponse	"Unauthorized - token is revoked, invalid, or Redis session expired"
// @Router			/v1/auth/refresh-token [post]
func (handler *AuthHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse request body
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest,
			response.SetErrorResponse(http.StatusBadRequest, "invalid request body"))
	}

	if req.RefreshToken == "" {
		return c.JSON(http.StatusUnauthorized,
			response.SetErrorResponse(http.StatusUnauthorized, constants.TokenRevokedMessage))
	}

	// Call usecase
	newToken, err := handler.AuthUseCase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == constants.ErrTokenRevoked {
			return c.JSON(http.StatusUnauthorized,
				response.SetErrorResponse(http.StatusUnauthorized, constants.TokenRevokedMessage))
		}
		return c.JSON(http.StatusUnauthorized,
			response.SetErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	// Return response
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(ResponseAuth{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	})

	return c.JSON(http.StatusOK, resp)
}
