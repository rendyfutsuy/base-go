package usecase

import (
	"time"

	auth "github.com/rendyfutsuy/base-go/modules/auth"
	roleManagement "github.com/rendyfutsuy/base-go/modules/role_management"
	"github.com/rendyfutsuy/base-go/modules/user_management"
	"github.com/rendyfutsuy/base-go/utils/services/queue"
)

type userUsecase struct {
	userRepo       user_management.Repository
	roleManagement roleManagement.Repository
	auth           auth.Repository
	queue          queue.QueueService
	contextTimeout time.Duration
}

func NewUserManagementUsecase(r user_management.Repository, rm roleManagement.Repository, auth auth.Repository, timeout time.Duration, q queue.QueueService) user_management.Usecase {
	return &userUsecase{
		userRepo:       r,
		roleManagement: rm,
		auth:           auth,
		queue:          q,
		contextTimeout: timeout,
	}
}
