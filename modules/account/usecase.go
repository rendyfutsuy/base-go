package account

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/account/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	CreateAccount(c echo.Context, req *dto.ReqCreateAccount, authId uuid.UUID) (accountRes *models.Account, err error)
	GetAccountByID(id string) (account *models.Account, err error)
	GetAllAccount() (accounts []models.Account, err error)
	GetIndexAccount(req request.PageRequest) (accounts []models.Account, total int, err error)
	UpdateAccount(id string, req *dto.ReqUpdateAccount, authId uuid.UUID) (accountRes *models.Account, err error)
	SoftDeleteAccount(id string, authId uuid.UUID) (accountRes *models.Account, err error)
}
