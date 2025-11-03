package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionGroupByID(c echo.Context, id string) (role *models.PermissionGroup, err error) {
	ctx := c.Request().Context()

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionGroupByID(ctx, uId)
}

func (u *roleUsecase) GetIndexPermissionGroup(c echo.Context, req request.PageRequest) (role_infos []models.PermissionGroup, total int, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetIndexPermissionGroup(ctx, req)
}

func (u *roleUsecase) GetAllPermissionGroup(c echo.Context) (role_infos []models.PermissionGroup, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetAllPermissionGroup(ctx)
}

func (u *roleUsecase) PermissionGroupNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (permissionGroupRes *models.PermissionGroup, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetDuplicatedPermissionGroup(ctx, name, id)
}
