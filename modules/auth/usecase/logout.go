package usecase

import (
	"context"

	"github.com/rendyfutsuy/base-go/utils"
)

func (u *authUsecase) SignOut(ctx context.Context, token string) error {

	// destroy requested jwt token
	err := u.authRepo.DestroyToken(ctx, token)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
