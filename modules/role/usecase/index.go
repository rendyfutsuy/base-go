package usecase

import (
	"time"

	role "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role"
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
