package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/utils"
)

func (u *authUsecase) SignOut(c echo.Context, token string) error {

	// destroy requested jwt token
	err := u.authRepo.DestroyToken(token)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
