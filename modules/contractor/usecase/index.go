package usecase

import (
	"time"

	contractor "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor"
)

type contractorUsecase struct {
	contractorRepo contractor.Repository
	contextTimeout time.Duration
}

func NewContractorUsecase(r contractor.Repository, timeout time.Duration) contractor.Usecase {
	return &contractorUsecase{
		contractorRepo: r,
		contextTimeout: timeout,
	}
}
