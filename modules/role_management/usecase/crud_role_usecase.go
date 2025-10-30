package usecase

import (
	"errors"
	"fmt"

	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (u *roleUsecase) CreateRole(c echo.Context, req *dto.ReqCreateRole, authId string) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()

	// assert each Permission group exists
	for _, permissionGroupId := range req.PermissionGroups {
		// check permission availability on DB
		_, err := u.roleRepo.GetPermissionGroupByID(ctx, permissionGroupId)

		// return error if any permission group not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Function with ID `%s` is not Found..", permissionGroupId))
		}
	}

	// assert name is not duplicated
	result, err := u.roleRepo.RoleNameIsNotDuplicated(ctx, req.Name, uuid.Nil)

	if err != nil {
		return nil, err
	}

	if result == false {
		utils.Logger.Error(constants.RoleErrorRoleNotFound)
		return nil, errors.New(constants.RoleErrorRoleNotFound)
	}

	count, err := u.roleRepo.CountRole(ctx)
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	roleDb := req.ToDBCreateRole(formatCount, authId)

	roleRes, err = u.roleRepo.CreateRole(ctx, roleDb)
	if err != nil {
		return nil, err
	}

	return roleRes, err
}

func (u *roleUsecase) GetRoleByID(c echo.Context, id string) (role *models.Role, err error) {
	ctx := c.Request().Context()

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetRoleByID(ctx, uId)
}

func (u *roleUsecase) GetIndexRole(c echo.Context, req request.PageRequest) (role_infos []models.Role, total int, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetIndexRole(ctx, req)
}

func (u *roleUsecase) GetAllRole(c echo.Context) (role_infos []models.Role, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetAllRole(ctx)
}

func (u *roleUsecase) UpdateRole(c echo.Context, id string, req *dto.ReqUpdateRole, authId string) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()

	// assert each Permission group exists
	for _, permissionGroupId := range req.PermissionGroups {
		// check permission availability on DB
		_, err := u.roleRepo.GetPermissionGroupByID(ctx, permissionGroupId)

		// return error if any permission group not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Function with ID `%s` is not Found..", permissionGroupId))
		}
	}

	// parsing UUID
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// assert name is not duplicated
	result, err := u.roleRepo.RoleNameIsNotDuplicated(ctx, req.Name, uId)

	if err != nil {
		return nil, err
	}

	if result == false {
		utils.Logger.Error(constants.RoleErrorRoleNotFound)
		return nil, errors.New(constants.RoleErrorRoleNotFound)
	}

	// Mapping Input to DB
	roleDb := dto.ToDBUpdateRole{
		Name:             req.Name,
		Description:      req.Description,
		Cobs:             req.Cobs,
		PermissionGroups: req.PermissionGroups,
		Categories:       req.Categories,
	}

	return u.roleRepo.UpdateRole(ctx, uId, roleDb)
}

func (u *roleUsecase) SoftDeleteRole(c echo.Context, id string, authId string) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()

	// if role has user, return error
	role, err := u.GetRoleByID(c, id)
	if err != nil {
		return nil, errors.New("Role Not Found")
	}

	if role.TotalUser > 0 {
		return nil, errors.New("Role has user. Can't be deleted")
	}

	roleDb := dto.ToDBDeleteRole{}

	return u.roleRepo.SoftDeleteRole(ctx, role.ID, roleDb)
}

func (u *roleUsecase) RoleNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetDuplicatedRole(ctx, name, id)
}

func (u *roleUsecase) MyPermissionsByUserToken(c echo.Context, token string) (role *models.Role, err error) {
	ctx := c.Request().Context()

	// get user id from token
	user, err := u.authRepo.GetUserByAccessToken(ctx, token)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}

	return u.roleRepo.GetRoleByID(ctx, user.RoleId)
}
