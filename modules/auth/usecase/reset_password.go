package usecase

import (
	"context"
	"errors"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) RequestResetPassword(ctx context.Context, email string) error {
	// get user by email
	_, err := u.authRepo.FindByEmailOrUsername(ctx, email)

	// if fail to get user return error
	if err != nil {
		return errors.New(constants.AuthEmailNotFound)

	}

	return u.authRepo.RequestResetPassword(ctx, email)
}

func (u *authUsecase) ResetUserPassword(ctx context.Context, newPassword string, token string) error {
	// find user by password token
	user, err := u.authRepo.GetUserByResetPasswordToken(ctx, token)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, newPassword, user.ID)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New(constants.AuthNewPasswordSameAsOld)
	}

	// assert new password not the same with any previous password
	isCurrentPasswordPassed, err := u.authRepo.AssertPasswordNeverUsesByUser(ctx, newPassword, user.ID)

	// if new password failed to passed return error
	if !isCurrentPasswordPassed {
		return err
	}

	// update user password bases on new_password
	_, err = u.authRepo.UpdatePasswordById(ctx, newPassword, user.ID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// update password expired at to 3 month from now
	err = u.authRepo.IncreasePasswordExpiredAt(ctx, user.ID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// add new password to password history
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// add new password to password history
	err = u.authRepo.AddPasswordHistory(ctx, string(hashedPassword), user.ID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// reset password attempt counter to 0
	err = u.authRepo.ResetPasswordAttempt(ctx, user.ID)

	// if fail to reset return error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// destroy all reset password token
	err = u.authRepo.DestroyAllResetPasswordToken(ctx, user.ID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// destroy all token session
	err = u.authRepo.DestroyAllToken(ctx, user.ID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
