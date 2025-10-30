package usecase

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"

	"github.com/google/uuid"
)

func (u *roleUsecase) GetPermissionByID(c echo.Context, id string) (role *models.Permission, err error) {
	ctx := c.Request().Context()

	uId, err := utils.StringToUUID(id)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}

	return u.roleRepo.GetPermissionByID(ctx, uId)
}

func (u *roleUsecase) GetIndexPermission(c echo.Context, req request.PageRequest) (role_infos []models.Permission, total int, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetIndexPermission(ctx, req)
}

func (u *roleUsecase) GetAllPermission(c echo.Context) (role_infos []models.Permission, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetAllPermission(ctx)
}

func (u *roleUsecase) PermissionNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (permissionRes *models.Permission, err error) {
	ctx := c.Request().Context()
	return u.roleRepo.GetDuplicatedPermission(ctx, name, id)
}
