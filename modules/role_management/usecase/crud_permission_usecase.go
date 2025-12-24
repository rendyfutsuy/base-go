package usecase

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionByID(ctx context.Context, id string) (role *models.Permission, err error) {
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionByID(ctx, uId)
}

func (u *roleUsecase) GetIndexPermission(ctx context.Context, req request.PageRequest) (role_infos []models.Permission, total int, err error) {
	return u.roleRepo.GetIndexPermission(ctx, req)
}

func (u *roleUsecase) GetAllPermission(ctx context.Context) (role_infos []models.Permission, err error) {
	return u.roleRepo.GetAllPermission(ctx)
}

func (u *roleUsecase) PermissionNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (permissionRes *models.Permission, err error) {
	return u.roleRepo.GetDuplicatedPermission(ctx, name, id)
}
