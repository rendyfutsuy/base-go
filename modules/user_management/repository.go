package user_management

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// ------------------------------------------------- user scope - BEGIN -----------------------------------------------------------
	// crud
	CreateUser(userReq dto.ToDBCreateUser) (userRes *models.User, err error)
	GetUserByID(id uuid.UUID) (user *models.User, err error)
	GetAllUser() (users []models.User, err error)
	GetIndexUser(req request.PageRequest, filter dto.ReqUserIndexFilter) (users []models.User, total int, err error)
	UpdateUser(id uuid.UUID, userReq dto.ToDBUpdateUser) (userRes *models.User, err error)
	SoftDeleteUser(id uuid.UUID, userReq dto.ToDBDeleteUser) (userRes *models.User, err error)
	UserNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedUser(name string, excludedId uuid.UUID) (user *models.User, err error)
	UserNameIsNotDuplicatedOnSoftDeleted(name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedUserOnSoftDeleted(name string, excludedId uuid.UUID) (user *models.User, err error)
	BlockUser(id uuid.UUID) (userRes *models.User, err error)
	UnBlockUser(id uuid.UUID) (userRes *models.User, err error)
	ActivateUser(id uuid.UUID) (userRes *models.User, err error)
	DisActivateUser(id uuid.UUID) (userRes *models.User, err error)
	EmailIsNotDuplicated(email string, excludedId uuid.UUID) (bool, error)

	CountUser() (count *int, err error)
	// ------------------------------------------------- user scope - END ----------------------------------------------------------

	// ------------------------------------------------- password scope - BEGIN -----------------------------------------------------
	IsUserPasswordCanUpdated(id uuid.UUID) (bool, error)
	// ------------------------------------------------- password scope - END -------------------------------------------------------
}
