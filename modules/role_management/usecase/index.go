package usecase

import (
	"time"

	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/role_management"
)

type roleUsecase struct {
	roleRepo       role_management.Repository
	authRepo       auth.Repository
	contextTimeout time.Duration
}

func NewRoleManagementUsecase(r role_management.Repository, a auth.Repository, timeout time.Duration) role_management.Usecase {
	return &roleUsecase{
		authRepo:       a,
		roleRepo:       r,
		contextTimeout: timeout,
	}
}
