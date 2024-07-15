package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

func (u *contractorUsecase) CreateContractor(c echo.Context, req *dto.ReqCreateContractor, authId string) (contractorRes *models.Contractor, err error) {

	count, err := u.contractorRepo.CountContractor()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	contractorDb := req.ToDBCreateContractor(formatCount, authId)

	contractorRes, err = u.contractorRepo.CreateContractor(contractorDb)
	if err != nil {
		return nil, err
	}

	return contractorRes, err
}

func (u *contractorUsecase) GetContractorByID(id string) (contractor *models.Contractor, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return u.contractorRepo.GetContractorByID(uId)
}

func (u *contractorUsecase) GetIndexContractor(req request.PageRequest) (contractores []models.Contractor, total int, err error) {
	return u.contractorRepo.GetIndexContractor(req)
}

func (u *contractorUsecase) GetAllContractor() (contractores []models.Contractor, err error) {

	return u.contractorRepo.GetAllContractor()
}

func (u *contractorUsecase) UpdateContractor(id string, req *dto.ReqUpdateContractor, authId string) (contractorRes *models.Contractor, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	contractorDb := req.ToDBUpdateContractor(authId)

	return u.contractorRepo.UpdateContractor(uId, contractorDb)
}

func (u *contractorUsecase) SoftDeleteContractor(id string, authId string) (contractorRes *models.Contractor, err error) {

	uId, err := utils.StringToUUID(id)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	contractorDb := dto.ToDBDeleteContractor{
		DeletedByID: authId,
	}

	return u.contractorRepo.SoftDeleteContractor(uId, contractorDb)
}
