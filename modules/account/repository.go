package account

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	models "github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/modules/account/dto"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// crud
	CreateAccount(accountReq dto.ToDBCreateAccount) (accountRes *models.Account, err error)
	GetAccountByID(id uuid.UUID) (account *models.Account, err error)
	GetIndexAccount(req request.PageRequest) (accounts []models.Account, total int, err error)
	GetAllAccount() (accounts []models.Account, err error)
	UpdateAccount(id uuid.UUID, accountReq dto.ToDBUpdateAccount) (accountRes *models.Account, err error)
	SoftDeleteAccount(id uuid.UUID, accountReq dto.ToDBDeleteAccount) (accountRes *models.Account, err error)

	// general
	CountAccount() (count *int, err error)
}
