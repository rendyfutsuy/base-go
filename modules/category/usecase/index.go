package usecase

import (
	"time"

	category "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category"
)

type categoryUsecase struct {
	categoryRepo   category.Repository
	contextTimeout time.Duration
}

func NewCategoryUsecase(r category.Repository, timeout time.Duration) category.Usecase {
	return &categoryUsecase{
		categoryRepo:   r,
		contextTimeout: timeout,
	}
}
