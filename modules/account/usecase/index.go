package usecase

import (
	"time"

	account "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/account"
)

type accountUsecase struct {
	accountRepo    account.Repository
	contextTimeout time.Duration
}

func NewAccountUsecase(r account.Repository, timeout time.Duration) account.Usecase {
	return &accountUsecase{
		accountRepo:    r,
		contextTimeout: timeout,
	}
}
