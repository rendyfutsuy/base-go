package usecase

import (
	"errors"
	"fmt"

	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

func (u *roleUsecase) ReAssignPermissionByGroup(roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (roleRes *models.Role, err error) {
	// assert each Permission group exists
	for _, permissionGroupId := range req.PermissionGroupIds {
		// check permission availability on DB
		_, err := u.roleRepo.GetPermissionGroupByID(permissionGroupId)

		// return error if any permission group not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Permission Group with ID `%s` is not Found..", permissionGroupId))
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
	err = u.roleRepo.ReAssignPermissionGroup(uId, permissionGroupDb)

	if err != nil {
		return nil, err
	}

	return u.roleRepo.GetRoleByID(uId)
}

func (u *roleUsecase) AssignUsersToRole(roleId string, req *dto.ReqUpdateAssignUsersToRole) (roleRes *models.Role, err error) {
	// assert each User exists
	for _, userId := range req.UserIds {
		// check user availability on DB
		_, err := u.roleRepo.GetUserByID(userId)

		// return error if any user not valid one.
		if err != nil {
			return nil, errors.New(fmt.Sprintf("User with ID `%s` is not Found..", userId))
		}
	}

	// parsing UUID
	uId, err := utils.StringToUUID(roleId)

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	// assert role exists
	_, err = u.roleRepo.GetRoleByID(uId)

	// return error if any role not valid one.
	if err != nil {
		return nil, errors.New("Role with ID `" + roleId + "` is not Found..")
	}

	// assign Users to role
	err = u.roleRepo.AssignUsers(uId, req.UserIds)

	if err != nil {
		return nil, errors.New("Something went wrong when assigning users to role, please check if role and users exist")
	}

	return u.roleRepo.GetRoleByID(uId)
}
