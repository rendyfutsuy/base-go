package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/utils"
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

	r.GET("/refresh-token",
		handler.RefreshToken,
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
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqAuthUser)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// check if user password already expired
	if err := handler.AuthUseCase.IsUserPasswordExpired(ctx, req.Login); err != nil {
		return c.JSON(419, GeneralResponse{Message: err.Error()})
	}

	// Assuming the Authenticate method on AuthUseCase does the actual authentication
	token, err := handler.AuthUseCase.Authenticate(ctx, req.Login, req.Password)
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
	ctx := c.Request().Context()

	// parse token
	token := c.Get("token").(string)

	// initiate session destroy
	err := handler.AuthUseCase.SignOut(ctx, token)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: constants.AuthLogoutSuccess})
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
	ctx := c.Request().Context()

	// parse token
	token := c.Get("token").(string)

	// initiate session destroy
	user, err := handler.AuthUseCase.GetProfile(ctx, token)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// Convert models.User to dto.UserProfile
	profile := dto.UserProfile{
		UserId: user.ID.String(),
		Name:   user.FullName,
		Email:  user.Email,
		Role: utils.NullString{
			String: user.RoleName,
			Valid:  user.RoleName != "",
		},
		Gender:           user.Gender,
		Permissions:      user.Permissions,
		PermissionGroups: user.PermissionGroups,
		Modules:          user.Modules,
	}
	if user.IsActive {
		profile.Status = "Active"
	} else {
		profile.Status = "In Active"
	}

	return c.JSON(http.StatusOK, profile)
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
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqUpdateProfile)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// parse user ID from context
	userId := c.Get("userId").(string)

	// call update profile function
	err := handler.AuthUseCase.UpdateProfile(ctx, *req, userId)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: constants.AuthProfileUpdated})
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
	ctx := c.Request().Context()

	// Validate input
	req := new(dto.ReqUpdatePassword)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// initiate validation
	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	// parse user ID from context
	userId := c.Get("userId").(string)

	// call update profile function
	err := handler.AuthUseCase.UpdateMyPassword(ctx, *req, userId)

	// return error, if something happen
	if err != nil {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, GeneralResponse{Message: constants.AuthPasswordUpdated})
}

// @Summary		Refresh access token
// @Description	Generates a new access token based on the provided bearer token. If the token is revoked, returns an error.
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	ResponseAuth					"Successfully refreshed token"
// @Failure		400	{object}	GeneralResponse{message=string}	"Token is revoked or invalid request"
// @Failure		401	{object}	GeneralResponse{message=string}	"Unauthorized"
// @Router			/v1/auth/refresh-token [get]
func (handler *AuthHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()

	// parse token from context (set by middleware)
	token, ok := c.Get("token").(string)
	if !ok || token == "" {
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: constants.TokenRevokedMessage})
	}

	// initiate refresh token
	newToken, err := handler.AuthUseCase.RefreshToken(ctx, token)
	if err != nil {
		// Check if error is about revoked token
		if err == constants.ErrTokenRevoked {
			return c.JSON(http.StatusBadRequest, GeneralResponse{Message: constants.TokenRevokedMessage})
		}
		return c.JSON(http.StatusBadRequest, GeneralResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseAuth{AccessToken: newToken})
}
