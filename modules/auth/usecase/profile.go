package usecase

import (
	"context"
	"errors"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	filedto "github.com/rendyfutsuy/base-go/modules/file/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/token_storage"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) GetProfile(ctx context.Context, accessToken string) (user models.User, err error) {
	user, err = token_storage.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return user, err
	}

	// Get permissions and permission groups from role if role_id exists
	permissions := []string{}
	permissionGroups := []string{}
	modules := []string{}
	moduleMap := make(map[string]bool) // Use map to track unique modules
	if user.RoleId != uuid.Nil {
		// Get permissions
		permissionList, err := u.roleManagementRepo.GetPermissionFromRoleId(ctx, user.RoleId)
		if err == nil && len(permissionList) > 0 {
			for _, permission := range permissionList {
				permissions = append(permissions, permission.Name)
			}
		}

		// Get permission groups and extract unique modules
		permissionGroupList, err := u.roleManagementRepo.GetPermissionGroupFromRoleId(ctx, user.RoleId)
		if err == nil && len(permissionGroupList) > 0 {
			for _, permissionGroup := range permissionGroupList {
				permissionGroups = append(permissionGroups, permissionGroup.Name)

				// Extract module if valid and not already in map
				if permissionGroup.Module.Valid && permissionGroup.Module.String != "" {
					moduleName := permissionGroup.Module.String
					if !moduleMap[moduleName] {
						modules = append(modules, moduleName)
						moduleMap[moduleName] = true
					}
				}
			}
		}
	}
	user.Permissions = permissions
	user.PermissionGroups = permissionGroups
	user.Modules = modules

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

	// Check if current user has already changed password (is_first_time_login = false)
	isFirstTimeLogin, err := u.authRepo.GetIsFirstTimeLogin(ctx, userUUID)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	if !isFirstTimeLogin {
		return errors.New(constants.AuthPasswordAlreadyChanged)
	}

	// assert current password not the same with new password
	isNewPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, passwordChunks.NewPassword, userUUID)

	// if current password same with current password, return error
	if isNewPasswordRight {
		return errors.New(constants.AuthNewPasswordSameAsOld)
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
	err = token_storage.RevokeAllUserSessions(ctx, userUUID)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (u *authUsecase) UpdateMyAvatar(ctx context.Context, user models.User, file *multipart.FileHeader) error {
	// parse file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// upload via File usecase
	extra := user.Username
	uploaded, err := u.fileUC.Upload(ctx, filedto.UploadInput{
		Data:             fileData,
		OriginalFileName: file.Filename,
		DestRoot:         "users/avatars",
		ExtraPath:        &extra,
	})
	if err != nil {
		return err
	}

	// ensure only one avatar: unassign old avatar records for this user
	_ = u.fileUC.UnassignFiles(ctx, filedto.UnassignFilesFromModule{
		ModuleID:   user.ID,
		ModuleType: constants.ModuleTypeUser,
	})

	// assign uploaded file as current avatar
	ftype := constants.FileTypeAvatar
	if err := u.fileUC.AssignFiles(ctx, filedto.AssignFilesToModule{
		ModuleID:   user.ID,
		ModuleType: constants.ModuleTypeUser,
		Items: []filedto.AssignFileItem{
			{FileID: uploaded.ID, Type: &ftype},
		},
	}); err != nil {
		return err
	}

	// no DB update needed; avatar served via pivot
	return nil
}

// UploadAvatar is deprecated in favor of file module integration
