package usecase

import (
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionByID(id string) (role *models.Permission, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionByID(uId)
}

func (u *roleUsecase) GetIndexPermission(req request.PageRequest) (role_infos []models.Permission, total int, err error) {
	return u.roleRepo.GetIndexPermission(req)
}

func (u *roleUsecase) GetAllPermission() (role_infos []models.Permission, err error) {

	return u.roleRepo.GetAllPermission()
}

func (u *roleUsecase) PermissionNameIsNotDuplicated(name string, id uuid.UUID) (permissionRes *models.Permission, err error) {
	return u.roleRepo.GetDuplicatedPermission(name, id)
}
