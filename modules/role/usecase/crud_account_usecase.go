package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (u *roleUsecase) CreateRole(c echo.Context, req *dto.ReqCreateRole, authId uuid.UUID) (roleRes *models.Role, err error) {

	count, err := u.roleRepo.CountRole()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	roleDb := dto.ToDBCreateRole{
		Name:        req.Name,
		Code:        formatCount,
		CreatedByID: authId,
	}

	roleRes, err = u.roleRepo.CreateRole(roleDb)
	if err != nil {
		return nil, err
	}

	return roleRes, err
}

func (u *roleUsecase) GetRoleByID(id string) (role *models.Role, err error) {

	uId := uuid.MustParse(id)

	return u.roleRepo.GetRoleByID(uId)
}

func (u *roleUsecase) GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error) {
	return u.roleRepo.GetIndexRole(req)
}

func (u *roleUsecase) GetAllRole() (roles []models.Role, err error) {

	return u.roleRepo.GetAllRole()
}

func (u *roleUsecase) UpdateRole(id string, req *dto.ReqUpdateRole, authId uuid.UUID) (roleRes *models.Role, err error) {

	uId := uuid.MustParse(id)

	roleDb := dto.ToDBUpdateRole{
		Name:        req.Name,
		UpdatedByID: authId,
	}

	return u.roleRepo.UpdateRole(uId, roleDb)
}

func (u *roleUsecase) SoftDeleteRole(id string, authId uuid.UUID) (roleRes *models.Role, err error) {

	uId := uuid.MustParse(id)

	roleDb := dto.ToDBDeleteRole{
		DeletedByID: authId,
	}

	return u.roleRepo.SoftDeleteRole(uId, roleDb)
}
