package user_management

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
)

type Usecase interface {
	// user scope
	CreateUser(ctx context.Context, req *dto.ReqCreateUser, authId string) (userRes *models.User, err error)
	GetUserByID(ctx context.Context, id string) (user *models.User, err error)
	GetAllUser(ctx context.Context) (user_infos []models.User, err error)
	GetIndexUser(ctx context.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (user_infos []models.User, total int, err error)
	UpdateUser(ctx context.Context, id string, req *dto.ReqUpdateUser, authId string) (userRes *models.User, err error)
	SoftDeleteUser(ctx context.Context, id string, authId string) (userRes *models.User, err error)
	UserNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (userRes *models.User, err error)
	EmailIsNotDuplicated(ctx context.Context, email string, id uuid.UUID) (userRes *models.User, err error)
	BlockUser(ctx context.Context, id string, req *dto.ReqBlockUser) (userRes *models.User, err error)
	ActivateUser(ctx context.Context, id string, req *dto.ReqActivateUser) (userRes *models.User, err error)

	// password management
	UpdateUserPassword(ctx context.Context, userId string, passwordChunks *dto.ReqUpdateUserPassword) error
	UpdateUserPasswordNoCheckRequired(ctx context.Context, userId string, passwordChunks *dto.ReqUpdateUserPassword) error
	AssertCurrentUserPassword(ctx context.Context, id string, inputtedPassword string) error

	// import users
	ImportUsersFromExcel(ctx context.Context, filePath string) (res *dto.ResImportUsers, err error)
}
