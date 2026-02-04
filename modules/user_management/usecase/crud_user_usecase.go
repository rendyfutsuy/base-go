package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/token_storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// validateRole checks if the role exists and is valid
func (u *userUsecase) validateRole(ctx context.Context, roleId uuid.UUID) error {
	if roleId == uuid.Nil {
		return nil
	}

	roleObject, err := u.roleManagement.GetRoleByID(ctx, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.UserRoleNotFound)
		}
		return err
	}

	if roleObject == nil || roleObject.ID == uuid.Nil {
		return errors.New(constants.UserRoleNotFound)
	}

	return nil
}

// validateUsernameNotDuplicated checks if username is not duplicated
func (u *userUsecase) validateUsernameNotDuplicated(ctx context.Context, username string, excludedId uuid.UUID) error {
	if username == "" {
		return nil
	}

	result, err := u.userRepo.UsernameIsNotDuplicated(ctx, username, excludedId)
	if err != nil {
		return err
	}

	if !result {
		utils.Logger.Error(constants.UserUsernameAlreadyExistsID)
		return errors.New(constants.UserUsernameAlreadyExistsID)
	}

	return nil
}

// validateEmailNotDuplicated checks if email is not duplicated
func (u *userUsecase) validateEmailNotDuplicated(ctx context.Context, email string, excludedId uuid.UUID) error {
	if email == "" {
		return nil
	}

	result, err := u.userRepo.EmailIsNotDuplicated(ctx, email, excludedId)
	if err != nil {
		return err
	}

	if !result {
		utils.Logger.Error(constants.UserEmailAlreadyExists)
		return errors.New(constants.UserEmailAlreadyExists)
	}

	return nil
}

func (u *userUsecase) CreateUser(ctx context.Context, req *dto.ReqCreateUser, userID string) (userRes *models.User, err error) {
	if err := u.validateRole(ctx, req.RoleId); err != nil {
		return nil, err
	}

	if err := u.validateUsernameNotDuplicated(ctx, req.Username, uuid.Nil); err != nil {
		return nil, err
	}

	count, err := u.userRepo.CountUser(ctx)
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)
	userDb := req.ToDBCreateUser(formatCount, userID)

	userRes, err = u.userRepo.CreateUser(ctx, userDb)
	if err != nil {
		return nil, err
	}

	return userRes, err
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (user *models.User, err error) {
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.userRepo.GetUserByID(ctx, uId)
}

func (u *userUsecase) GetIndexUser(ctx context.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (user_infos []models.User, total int, err error) {
	return u.userRepo.GetIndexUser(ctx, req, filter)
}

func (u *userUsecase) GetAllUser(ctx context.Context) (user_infos []models.User, err error) {
	return u.userRepo.GetAllUser(ctx)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id string, req *dto.ReqUpdateUser, userID string) (userRes *models.User, err error) {
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	if err := u.validateRole(ctx, req.RoleId); err != nil {
		return nil, err
	}

	if err := u.validateUsernameNotDuplicated(ctx, req.Username, uId); err != nil {
		return nil, err
	}

	if err := u.validateEmailNotDuplicated(ctx, req.Email, uId); err != nil {
		return nil, err
	}

	userDb := dto.ToDBUpdateUser{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		IsActive: req.IsActive,
		RoleId:   req.RoleId,
		Gender:   req.Gender,
	}

	return u.userRepo.UpdateUser(ctx, uId, userDb)
}

func (u *userUsecase) SoftDeleteUser(ctx context.Context, id string, userID string) (userRes *models.User, err error) {
	// Check if user exists
	user, err := u.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}

	// Check if user is deletable
	if !user.Deletable {
		return nil, errors.New(constants.UserCannotDelete)
	}

	// TBA there would be more
	// but the other condition would integrate in another occasion

	userDb := dto.ToDBDeleteUser{}

	return u.userRepo.SoftDeleteUser(ctx, user.ID, userDb)
}

func (u *userUsecase) UserNameIsNotDuplicated(ctx context.Context, fullName string, id uuid.UUID) (userRes *models.User, err error) {
	return u.userRepo.GetDuplicatedUser(ctx, fullName, id)
}

func (u *userUsecase) BlockUser(ctx context.Context, id string, req *dto.ReqBlockUser) (userRes *models.User, err error) {
	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// determinate if user is block or not
	if !req.IsBlock {
		// user requested to be unblock
		// unblock user
		_, err = u.userRepo.UnBlockUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	} else {
		// user requested to be block
		// block user
		_, err = u.userRepo.BlockUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	}

	// revoke user token
	token_storage.RevokeAllUserSessions(ctx, uId)

	return u.userRepo.GetUserByID(ctx, uId)
}

func (u *userUsecase) ActivateUser(ctx context.Context, id string, req *dto.ReqActivateUser) (userRes *models.User, err error) {
	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// determinate if user is block or not
	if !req.IsActive {
		// user requested to be dis-activate
		// dis-activate user
		_, err = u.userRepo.DisActivateUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	} else {
		// user requested to be active
		// active user
		_, err = u.userRepo.ActivateUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	}

	// revoke user token
	token_storage.RevokeAllUserSessions(ctx, uId)

	return u.userRepo.GetUserByID(ctx, uId)
}
