package usecase

import (
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionGroupByID(id string) (role *models.PermissionGroup, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionGroupByID(uId)
}

func (u *roleUsecase) GetIndexPermissionGroup(req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error) {
	return u.roleRepo.GetIndexPermissionGroup(req)
}

func (u *roleUsecase) GetAllPermissionGroup() (role_infos []models.PermissionGroup, err error) {

	return u.roleRepo.GetAllPermissionGroup()
}

func (u *roleUsecase) PermissionGroupNameIsNotDuplicated(name string, id uuid.UUID) (permissionGroupRes *models.PermissionGroup, err error) {
	return u.roleRepo.GetDuplicatedPermissionGroup(name, id)
}
