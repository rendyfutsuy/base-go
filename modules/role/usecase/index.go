package usecase

import (
	"time"

	role "github.com/rendyfutsuy/base-go/modules/role"
)

type roleUsecase struct {
	roleRepo       role.Repository
	contextTimeout time.Duration
}

func NewRoleUsecase(r role.Repository, timeout time.Duration) role.Usecase {
	return &roleUsecase{
		roleRepo:       r,
		contextTimeout: timeout,
	}
}
