package usecase

import (
	"errors"

	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *userUsecase) UpdateUserPassword(id string, passwordChunks *dto.ReqUpdateUserPassword) error {
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
	isPasswordRight, err := u.auth.AssertPasswordRight(passwordChunks.OldPassword, userId)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Old Password not Match")
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.auth.AssertPasswordRight(passwordChunks.NewPassword, userId)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New("New Password should not be same with Current Password")
	}

	// assert new password not the same wit any previous password
	isCurrentPasswordPassed, err := u.auth.AssertPasswordNeverUsesByUser(passwordChunks.NewPassword, userId)

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
	err = u.auth.AddPasswordHistory(string(hashedPassword), userId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// reset password attempt counter to 0
	err = u.auth.ResetPasswordAttempt(userId)

	// if fail to reset return error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// update user password bases on new_password
	_, err = u.auth.UpdatePasswordById(passwordChunks.NewPassword, userId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// destroy all token session
	err = u.auth.DestroyAllToken(userId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *userUsecase) AssertCurrentUserPassword(id string, inputtedPassword string) error {
	// parsing UUID
	userId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert current password given is same with saved password
	isPasswordRight, err := u.auth.AssertPasswordRight(inputtedPassword, userId)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Given Password not Match with Current Password")
	}

	return nil
}
