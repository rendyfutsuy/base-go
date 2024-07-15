package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func (u *classUsecase) CreateClass(c echo.Context, req *dto.ReqCreateClass, authId string) (classRes *models.Class, err error) {

	count, err := u.classRepo.CountClass()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	classDb := req.ToDBCreateClass(formatCount, authId)

	classRes, err = u.classRepo.CreateClass(classDb)
	if err != nil {
		return nil, err
	}

	return classRes, err
}

func (u *classUsecase) GetClassByID(id string) (class *models.Class, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return u.classRepo.GetClassByID(uId)
}

func (u *classUsecase) GetIndexClass(req request.PageRequest) (classes []models.Class, total int, err error) {
	return u.classRepo.GetIndexClass(req)
}

func (u *classUsecase) GetAllClass() (classes []models.Class, err error) {

	return u.classRepo.GetAllClass()
}

func (u *classUsecase) UpdateClass(id string, req *dto.ReqUpdateClass, authId string) (classRes *models.Class, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	classDb := dto.ToDBUpdateClass{
		Name:        req.Name,
		UpdatedByID: authId,
	}

	return u.classRepo.UpdateClass(uId, classDb)
}

func (u *classUsecase) SoftDeleteClass(id string, authId string) (classRes *models.Class, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	classDb := dto.ToDBDeleteClass{
		DeletedByID: authId,
	}

	return u.classRepo.SoftDeleteClass(uId, classDb)
}
