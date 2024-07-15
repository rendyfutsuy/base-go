package usecase

import (
	"errors"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) RequestResetPassword(c echo.Context, email string) error {
	// get user by email
	_, err := u.authRepo.FindByEmailOrUsername(email)

	// if fail to get user return error
	if err != nil {
		return errors.New("Email not found, please be sure no typo on email input...")

	}

	return u.authRepo.RequestResetPassword(email)
}

func (u *authUsecase) ResetUserPassword(c echo.Context, newPassword string, token string) error {
	// find user by password token
	user, err := u.authRepo.GetUserByResetPasswordToken(token)

	if err != nil {
		return err
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.authRepo.AssertPasswordRight(newPassword, user.ID)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New("New Password should not be same with Current Password")
	}

	// assert new password not the same with any previous password
	isCurrentPasswordPassed, err := u.authRepo.AssertPasswordNeverUsesByUser(newPassword, user.ID)

	// if new password failed to passed return error
	if !isCurrentPasswordPassed {
		return err
	}

	// update user password bases on new_password
	_, err = u.authRepo.UpdatePasswordById(newPassword, user.ID)

	if err != nil {
		return err
	}

	// update password expired at to 3 month from now
	err = u.authRepo.IncreasePasswordExpiredAt(user.ID)

	if err != nil {
		return err
	}

	// add new password to password history
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// add new password to password history
	err = u.authRepo.AddPasswordHistory(string(hashedPassword), user.ID)

	if err != nil {
		return err
	}

	// destroy all reset password token
	err = u.authRepo.DestroyAllResetPasswordToken(user.ID)

	if err != nil {
		return err
	}

	// destroy all token session
	err = u.authRepo.DestroyAllToken(user.ID)

	if err != nil {
		return err
	}

	return nil
}
