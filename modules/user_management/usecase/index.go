package usecase

import (
	"time"

	"github.com/hibiken/asynq"
	auth "github.com/rendyfutsuy/base-go/modules/auth"
	roleManagement "github.com/rendyfutsuy/base-go/modules/role_management"
	"github.com/rendyfutsuy/base-go/modules/user_management"
)

type userUsecase struct {
	userRepo       user_management.Repository
	roleManagement roleManagement.Repository
	auth           auth.Repository
	queueClient    *asynq.Client
	contextTimeout time.Duration
}

func NewUserManagementUsecase(r user_management.Repository, rm roleManagement.Repository, auth auth.Repository, timeout time.Duration, queueClient *asynq.Client) user_management.Usecase {
	return &userUsecase{
		userRepo:       r,
		roleManagement: rm,
		auth:           auth,
		queueClient:    queueClient,
		contextTimeout: timeout,
	}
}
