package usecase

import (
	"time"

	class "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class"
)

type classUsecase struct {
	classRepo      class.Repository
	contextTimeout time.Duration
}

func NewClassUsecase(r class.Repository, timeout time.Duration) class.Usecase {
	return &classUsecase{
		classRepo:      r,
		contextTimeout: timeout,
	}
}
