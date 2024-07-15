package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func (u *carriageUsecase) CreateCarriage(c echo.Context, req *dto.ReqCreateCarriage, authId string) (carriageRes *models.Carriage, err error) {

	count, err := u.carriageRepo.CountCarriage()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	carriageDb := req.ToDBCreateCarriage(formatCount, authId)

	carriageRes, err = u.carriageRepo.CreateCarriage(carriageDb)
	if err != nil {
		return nil, err
	}

	return carriageRes, err
}

func (u *carriageUsecase) GetCarriageByID(id string) (carriage *models.Carriage, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return u.carriageRepo.GetCarriageByID(uId)
}

func (u *carriageUsecase) GetIndexCarriage(req request.PageRequest) (carriagees []models.Carriage, total int, err error) {
	return u.carriageRepo.GetIndexCarriage(req)
}

func (u *carriageUsecase) GetAllCarriage() (carriagees []models.Carriage, err error) {

	return u.carriageRepo.GetAllCarriage()
}

func (u *carriageUsecase) UpdateCarriage(id string, req *dto.ReqUpdateCarriage, authId string) (carriageRes *models.Carriage, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	carriageDb := req.ToDBUpdateCarriage(authId)

	return u.carriageRepo.UpdateCarriage(uId, carriageDb)
}

func (u *carriageUsecase) SoftDeleteCarriage(id string, authId string) (carriageRes *models.Carriage, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	carriageDb := dto.ToDBDeleteCarriage{
		DeletedByID: authId,
	}

	return u.carriageRepo.SoftDeleteCarriage(uId, carriageDb)
}
