package usecase

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

func (u *roleUsecase) ReAssignPermissionByGroup(c echo.Context, roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()

	// assert each Permission group exists
	for _, permissionGroupId := range req.PermissionGroupIds {
		// check permission availability on DB
		_, err := u.roleRepo.GetPermissionGroupByID(ctx, permissionGroupId)

		// return error if any permission group not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf(constants.PermissionGroupNotFoundWithIDAlt, permissionGroupId))
		}
	}

	// parsing UUID
	uId, err := utils.StringToUUID(roleId)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// Mapping Input to DB
	permissionGroupDb := dto.ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: req.PermissionGroupIds,
	}

	// re-assign permission groups to role
	err = u.roleRepo.ReAssignPermissionGroup(ctx, uId, permissionGroupDb)

	if err != nil {
		return nil, err
	}

	return u.roleRepo.GetRoleByID(ctx, uId)
}

func (u *roleUsecase) AssignUsersToRole(c echo.Context, roleId string, req *dto.ReqUpdateAssignUsersToRole) (roleRes *models.Role, err error) {
	ctx := c.Request().Context()

	// assert each User exists
	for _, userId := range req.UserIds {
		// check user availability on DB
		_, err := u.roleRepo.GetUserByID(ctx, userId)

		// return error if any user not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf(constants.UserNotFoundWithID, userId))
		}
	}

	// parsing UUID
	uId, err := utils.StringToUUID(roleId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// assert role exists
	_, err = u.roleRepo.GetRoleByID(ctx, uId)

	// return error if any role not valid one.
	if err != nil {
		return nil, errors.New(fmt.Sprintf(constants.RoleNotFoundWithID, roleId))
	}

	// assign Users to role
	err = u.roleRepo.AssignUsers(ctx, uId, req.UserIds)

	if err != nil {
		return nil, errors.New(constants.RoleAssignUsersError)
	}

	return u.roleRepo.GetRoleByID(ctx, uId)
}
