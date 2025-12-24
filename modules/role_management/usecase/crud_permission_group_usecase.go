package usecase

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionGroupByID(ctx context.Context, id string) (role *models.PermissionGroup, err error) {
	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionGroupByID(ctx, uId)
}

func (u *roleUsecase) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error) {
	return u.roleRepo.GetIndexPermissionGroup(ctx, req)
}

func (u *roleUsecase) GetAllPermissionGroup(ctx context.Context) (role_infos []models.PermissionGroup, err error) {
	return u.roleRepo.GetAllPermissionGroup(ctx)
}

func (u *roleUsecase) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (permissionGroupRes *models.PermissionGroup, err error) {
	return u.roleRepo.GetDuplicatedPermissionGroup(ctx, name, id)
}
