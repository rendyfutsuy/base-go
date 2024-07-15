package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func (u *conveyanceUsecase) CreateConveyance(c echo.Context, req *dto.ReqCreateConveyance, authId string) (conveyanceRes *models.Conveyance, err error) {

	count, err := u.conveyanceRepo.CountConveyance()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	conveyanceDb := req.ToDBCreateConveyance(formatCount, authId)

	conveyanceRes, err = u.conveyanceRepo.CreateConveyance(conveyanceDb)
	if err != nil {
		return nil, err
	}

	return conveyanceRes, err
}

func (u *conveyanceUsecase) GetConveyanceByID(id string) (conveyance *models.Conveyance, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return u.conveyanceRepo.GetConveyanceByID(uId)
}

func (u *conveyanceUsecase) GetIndexConveyance(req request.PageRequest) (conveyancees []models.Conveyance, total int, err error) {
	return u.conveyanceRepo.GetIndexConveyance(req)
}

func (u *conveyanceUsecase) GetAllConveyance() (conveyancees []models.Conveyance, err error) {

	return u.conveyanceRepo.GetAllConveyance()
}

func (u *conveyanceUsecase) UpdateConveyance(id string, req *dto.ReqUpdateConveyance, authId string) (conveyanceRes *models.Conveyance, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	conveyanceDb := req.ToDBUpdateConveyance(authId)

	return u.conveyanceRepo.UpdateConveyance(uId, conveyanceDb)
}

func (u *conveyanceUsecase) SoftDeleteConveyance(id string, authId string) (conveyanceRes *models.Conveyance, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	conveyanceDb := dto.ToDBDeleteConveyance{
		DeletedByID: authId,
	}

	return u.conveyanceRepo.SoftDeleteConveyance(uId, conveyanceDb)
}
