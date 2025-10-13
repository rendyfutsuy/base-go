package usecase

import (
	"fmt"

	"github.com/rendyfutsuy/base-go/helper/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/account/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (u *accountUsecase) CreateAccount(c echo.Context, req *dto.ReqCreateAccount, authId uuid.UUID) (accountRes *models.Account, err error) {

	count, err := u.accountRepo.CountAccount()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	accountDb := dto.ToDBCreateAccount{
		Name:        req.Name,
		Code:        formatCount,
		CreatedByID: authId,
	}

	accountRes, err = u.accountRepo.CreateAccount(accountDb)
	if err != nil {
		return nil, err
	}

	return accountRes, err
}

func (u *accountUsecase) GetAccountByID(id string) (account *models.Account, err error) {

	uId := uuid.MustParse(id)

	return u.accountRepo.GetAccountByID(uId)
}

func (u *accountUsecase) GetIndexAccount(req request.PageRequest) (accounts []models.Account, total int, err error) {
	return u.accountRepo.GetIndexAccount(req)
}

func (u *accountUsecase) GetAllAccount() (accounts []models.Account, err error) {

	return u.accountRepo.GetAllAccount()
}

func (u *accountUsecase) UpdateAccount(id string, req *dto.ReqUpdateAccount, authId uuid.UUID) (accountRes *models.Account, err error) {

	uId := uuid.MustParse(id)

	accountDb := dto.ToDBUpdateAccount{
		Name:        req.Name,
		UpdatedByID: authId,
	}

	return u.accountRepo.UpdateAccount(uId, accountDb)
}

func (u *accountUsecase) SoftDeleteAccount(id string, authId uuid.UUID) (accountRes *models.Account, err error) {

	uId := uuid.MustParse(id)

	accountDb := dto.ToDBDeleteAccount{
		DeletedByID: authId,
	}

	return u.accountRepo.SoftDeleteAccount(uId, accountDb)
}
