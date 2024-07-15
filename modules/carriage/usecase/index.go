package usecase

import (
	"time"

	carriage "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage"
)

type carriageUsecase struct {
	carriageRepo   carriage.Repository
	contextTimeout time.Duration
}

func NewCarriageUsecase(r carriage.Repository, timeout time.Duration) carriage.Usecase {
	return &carriageUsecase{
		carriageRepo:   r,
		contextTimeout: timeout,
	}
}
