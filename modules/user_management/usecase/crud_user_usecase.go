package usecase

import (
	"errors"
	"fmt"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (u *userUsecase) CreateUser(c echo.Context, req *dto.ReqCreateUser, authId string) (userRes *models.User, err error) {
	ctx := c.Request().Context()

	// assert email is not duplicated
	result, err := u.userRepo.EmailIsNotDuplicated(ctx, req.Email, uuid.Nil)

	if err != nil {
		return nil, err
	}

	if result == false {
		utils.Logger.Error(constants.UserEmailAlreadyExists)
		return nil, errors.New(constants.UserEmailAlreadyExists)
	}

	count, err := u.userRepo.CountUser(ctx)
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	userDb := req.ToDBCreateUser(formatCount, authId)

	userRes, err = u.userRepo.CreateUser(ctx, userDb)
	if err != nil {
		return nil, err
	}

	// Create New Password
	// update password
	// update user password bases on new_password
	_, err = u.auth.UpdatePasswordById(ctx, req.Password, userRes.ID)

	if err != nil {
		return userRes, err
	}

	return userRes, err
}

func (u *userUsecase) GetUserByID(c echo.Context, id string) (user *models.User, err error) {
	ctx := c.Request().Context()

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.userRepo.GetUserByID(ctx, uId)
}

func (u *userUsecase) GetIndexUser(c echo.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (user_infos []models.User, total int, err error) {
	ctx := c.Request().Context()
	return u.userRepo.GetIndexUser(ctx, req, filter)
}

func (u *userUsecase) GetAllUser(c echo.Context) (user_infos []models.User, err error) {
	ctx := c.Request().Context()
	return u.userRepo.GetAllUser(ctx)
}

func (u *userUsecase) UpdateUser(c echo.Context, id string, req *dto.ReqUpdateUser, authId string) (userRes *models.User, err error) {
	ctx := c.Request().Context()

	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// assert email is not duplicated
	result, err := u.userRepo.EmailIsNotDuplicated(ctx, req.Email, uId)

	if err != nil {
		return nil, err
	}

	if result == false {
		utils.Logger.Error(constants.UserEmailAlreadyExists)
		return nil, errors.New(constants.UserEmailAlreadyExists)
	}

	// Mapping Input to DB
	userDb := dto.ToDBUpdateUser{
		FullName: req.FullName,
		Email:    req.Email,
		IsActive: req.IsActive,
		RoleId:   req.RoleId,
		Gender:   req.Gender,
	}

	return u.userRepo.UpdateUser(ctx, uId, userDb)
}

func (u *userUsecase) SoftDeleteUser(c echo.Context, id string, authId string) (userRes *models.User, err error) {
	ctx := c.Request().Context()

	// if user has user, return error
	user, err := u.GetUserByID(c, id)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}

	userDb := dto.ToDBDeleteUser{}

	return u.userRepo.SoftDeleteUser(ctx, user.ID, userDb)
}

func (u *userUsecase) UserNameIsNotDuplicated(c echo.Context, fullName string, id uuid.UUID) (userRes *models.User, err error) {
	ctx := c.Request().Context()
	return u.userRepo.GetDuplicatedUser(ctx, fullName, id)
}

func (u *userUsecase) BlockUser(c echo.Context, id string, req *dto.ReqBlockUser) (userRes *models.User, err error) {
	ctx := c.Request().Context()

	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// determinate if user is block or not
	if req.IsBlock == false {
		// user requested to be unblock
		// unblock user
		_, err = u.userRepo.UnBlockUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	} else if req.IsBlock == true {
		// user requested to be block
		// block user
		_, err = u.userRepo.BlockUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	}

	// revoke user token
	u.auth.DestroyAllToken(ctx, uId)

	return u.userRepo.GetUserByID(ctx, uId)
}

func (u *userUsecase) ActivateUser(c echo.Context, id string, req *dto.ReqActivateUser) (userRes *models.User, err error) {
	ctx := c.Request().Context()

	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// determinate if user is block or not
	if req.IsActive == false {
		// user requested to be dis-activate
		// dis-activate user
		_, err = u.userRepo.DisActivateUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	} else if req.IsActive == true {
		// user requested to be active
		// active user
		_, err = u.userRepo.ActivateUser(ctx, uId)
		if err != nil {
			return nil, err
		}
	}

	// revoke user token
	u.auth.DestroyAllToken(ctx, uId)

	return u.userRepo.GetUserByID(ctx, uId)
}
