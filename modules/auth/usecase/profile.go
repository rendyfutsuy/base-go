package usecase

import (
	"errors"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/auth/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) GetProfile(c echo.Context) (profile dto.UserProfile, err error) {
	user, err := u.authRepo.FindByCurrentSession(c.Get("token").(string))
	return user, nil
}

func (u *authUsecase) UpdateProfile(c echo.Context, profileChunks dto.ReqUpdateProfile) error {
	// find user by token
	user, err := u.authRepo.FindByCurrentSession(c.Get("token").(string))

	if err != nil {
		return err
	}

	// parse user ID to UUID
	userId, err := uuid.Parse(user.UserId)
	if err != nil {
		return err
	}

	// update user profile
	// column updated: name
	_, err = u.authRepo.UpdateProfileById(profileChunks, userId)

	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) UpdateMyPassword(c echo.Context, passwordChunks dto.ReqUpdatePassword) error {
	// find user by token
	user, err := u.authRepo.FindByCurrentSession(c.Get("token").(string))

	if err != nil {
		return err
	}

	// parse user ID to UUID
	userId, err := uuid.Parse(user.UserId)
	if err != nil {
		return err
	}

	// assert old password given is same with saved password
	isPasswordRight, err := u.authRepo.AssertPasswordRight(passwordChunks.OldPassword, userId)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Old Password not Match")
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.authRepo.AssertPasswordRight(passwordChunks.NewPassword, userId)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New("New Password should not be same with Current Password")
	}

	// assert new password not the same wit any previous password
	isCurrentPasswordPassed, err := u.authRepo.AssertPasswordNeverUsesByUser(passwordChunks.NewPassword, userId)

	// if new password fail to match return error
	if !isCurrentPasswordPassed {
		return err
	}

	// add new password to password history
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChunks.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// add new password to password history
	err = u.authRepo.AddPasswordHistory(string(hashedPassword), userId)

	if err != nil {
		return err
	}

	// update user password bases on new_password
	_, err = u.authRepo.UpdatePasswordById(passwordChunks.NewPassword, userId)

	if err != nil {
		return err
	}

	// destroy all token session
	err = u.authRepo.DestroyAllToken(userId)
	if err != nil {
		return err
	}

	return nil
}
