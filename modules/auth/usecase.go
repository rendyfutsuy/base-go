package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// Usecase represent the auth's usecases
type Usecase interface {
	// every new usecase on ..modules/auth/usecase/, please register it here
	Authenticate(c echo.Context, login string, password string) (string, error)
	SignOut(c echo.Context, token string) error
	GetProfile(c echo.Context) (profile dto.UserProfile, err error)
	UpdateProfile(c echo.Context, profileChunks dto.ReqUpdateProfile) error
	UpdateMyPassword(c echo.Context, passwordChunks dto.ReqUpdatePassword) error

	// for reset password
	RequestResetPassword(c echo.Context, email string) error
	ResetUserPassword(c echo.Context, newPassword string, token string) error
}
