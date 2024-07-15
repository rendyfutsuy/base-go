package usecase

import (
	"time"

	conveyance "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance"
)

type conveyanceUsecase struct {
	conveyanceRepo conveyance.Repository
	contextTimeout time.Duration
}

func NewConveyanceUsecase(r conveyance.Repository, timeout time.Duration) conveyance.Usecase {
	return &conveyanceUsecase{
		conveyanceRepo: r,
		contextTimeout: timeout,
	}
}
