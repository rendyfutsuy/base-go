package usecase

import (
	"github.com/labstack/echo/v4"
)

func (u *authUsecase) SignOut(c echo.Context, token string) error {

	// destroy requested jwt token
	err := u.authRepo.DestroyToken(token)

	if err != nil {
		return err
	}

	return nil
}
