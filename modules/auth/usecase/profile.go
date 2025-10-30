package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) GetProfile(ctx context.Context, accessToken string) (profile dto.UserProfile, err error) {
	user, err := u.authRepo.FindByCurrentSession(ctx, accessToken)
	if err != nil {
		return profile, err
	}

	return user, nil
}

func (u *authUsecase) UpdateProfile(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId string) error {
	// parse user ID to UUID
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// update user profile
	// column updated: name
	_, err = u.authRepo.UpdateProfileById(ctx, profileChunks, userUUID)

	if err != nil {
		// utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *authUsecase) UpdateMyPassword(ctx context.Context, passwordChunks dto.ReqUpdatePassword, userId string) error {
	// parse user ID to UUID
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// assert old password given is same with saved password
	isPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, passwordChunks.OldPassword, userUUID)

	// if old password fail to match return error
	if !isPasswordRight {
		return errors.New("Old Password not Match")
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, passwordChunks.NewPassword, userUUID)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New("New Password should not be same with Current Password")
	}

	// assert new password not the same wit any previous password
	isCurrentPasswordPassed, err := u.authRepo.AssertPasswordNeverUsesByUser(ctx, passwordChunks.NewPassword, userUUID)

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
	err = u.authRepo.AddPasswordHistory(ctx, string(hashedPassword), userUUID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// reset password attempt counter to 0
	err = u.authRepo.ResetPasswordAttempt(ctx, userUUID)

	// if fail to reset return error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// update user password bases on new_password
	_, err = u.authRepo.UpdatePasswordById(ctx, passwordChunks.NewPassword, userUUID)

	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// destroy all token session
	err = u.authRepo.DestroyAllToken(ctx, userUUID)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}
