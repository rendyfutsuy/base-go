package user_management

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
)

type Usecase interface {
	// user scope
	CreateUser(c echo.Context, req *dto.ReqCreateUser, authId string) (userRes *models.User, err error)
	GetUserByID(c echo.Context, id string) (user *models.User, err error)
	GetAllUser(c echo.Context) (user_infos []models.User, err error)
	GetIndexUser(c echo.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (user_infos []models.User, total int, err error)
	UpdateUser(c echo.Context, id string, req *dto.ReqUpdateUser, authId string) (userRes *models.User, err error)
	SoftDeleteUser(c echo.Context, id string, authId string) (userRes *models.User, err error)
	UserNameIsNotDuplicated(c echo.Context, name string, id uuid.UUID) (userRes *models.User, err error)
	BlockUser(c echo.Context, id string, req *dto.ReqBlockUser) (userRes *models.User, err error)
	ActivateUser(c echo.Context, id string, req *dto.ReqActivateUser) (userRes *models.User, err error)

	// password management
	UpdateUserPassword(c echo.Context, userId string, passwordChunks *dto.ReqUpdateUserPassword) error
	AssertCurrentUserPassword(c echo.Context, id string, inputtedPassword string) error
}
