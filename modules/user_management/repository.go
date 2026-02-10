package user_management

import (
	"context"

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
	CreateUser(ctx context.Context, userReq dto.ToDBCreateUser) (userRes *models.User, err error)
	GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error)
	GetAllUser(ctx context.Context) (users []models.User, err error)
	GetIndexUser(ctx context.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (users []models.User, total int, err error)
	UpdateUser(ctx context.Context, id uuid.UUID, userReq dto.ToDBUpdateUser) (userRes *models.User, err error)
	SoftDeleteUser(ctx context.Context, id uuid.UUID, userReq dto.ToDBDeleteUser) (userRes *models.User, err error)
	UserNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedUser(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error)
	GetDuplicatedUserByEmail(ctx context.Context, email string, excludedId uuid.UUID) (user *models.User, err error)
	UserNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error)
	GetDuplicatedUserOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error)
	BlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error)
	UnBlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error)
	ActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error)
	DisActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error)
	EmailIsNotDuplicated(ctx context.Context, email string, excludedId uuid.UUID) (bool, error)
	UsernameIsNotDuplicated(ctx context.Context, username string, excludedId uuid.UUID) (bool, error)
	NikIsNotDuplicated(ctx context.Context, nik string, excludedId uuid.UUID) (bool, error)
	CheckBatchDuplication(ctx context.Context, emails, usernames, niks []string) (duplicatedEmails, duplicatedUsernames, duplicatedNiks map[string]bool, err error)
	BulkCreateUsers(ctx context.Context, usersReq []dto.ToDBCreateUser) (err error)

	CountUser(ctx context.Context) (count *int, err error)
	// ------------------------------------------------- user scope - END ----------------------------------------------------------

	// ------------------------------------------------- password scope - BEGIN -----------------------------------------------------
	IsUserPasswordCanUpdated(ctx context.Context, id uuid.UUID) (bool, error)
	// ------------------------------------------------- password scope - END -------------------------------------------------------
}
