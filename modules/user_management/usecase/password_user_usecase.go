package usecase

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUsecase) UpdateUserPassword(c echo.Context, id string, passwordChunks *dto.ReqUpdateUserPassword) error {
	ctx := c.Request().Context()

	// parsing UUID
	userId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert user password can be change
	_, err = u.userRepo.IsUserPasswordCanUpdated(userId)

	// if error occurs, return error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert old password given is same with saved password
	isPasswordRight, err := u.auth.AssertPasswordRight(ctx, passwordChunks.OldPassword, userId)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Old Password not Match")
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.auth.AssertPasswordRight(ctx, passwordChunks.NewPassword, userId)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New("New Password should not be same with Current Password")
	}

	// assert new password not the same wit any previous password
	isCurrentPasswordPassed, err := u.auth.AssertPasswordNeverUsesByUser(ctx, passwordChunks.NewPassword, userId)

	// if new password fail to match return error
	if !isCurrentPasswordPassed {
		return err
	}

	// add new password to password history
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChunks.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// add new password to password history
	err = u.auth.AddPasswordHistory(ctx, string(hashedPassword), userId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// reset password attempt counter to 0
	err = u.auth.ResetPasswordAttempt(ctx, userId)

	// if fail to reset return error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// update user password bases on new_password
	_, err = u.auth.UpdatePasswordById(ctx, passwordChunks.NewPassword, userId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// destroy all token session
	err = u.auth.DestroyAllToken(ctx, userId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *userUsecase) AssertCurrentUserPassword(c echo.Context, id string, inputtedPassword string) error {
	ctx := c.Request().Context()

	// parsing UUID
	userId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert current password given is same with saved password
	isPasswordRight, err := u.auth.AssertPasswordRight(ctx, inputtedPassword, userId)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Given Password not Match with Current Password")
	}

	return nil
}
