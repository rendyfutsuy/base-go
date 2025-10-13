package usecase

import (
	"time"

	account "github.com/rendyfutsuy/base-go.git/modules/account"
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
