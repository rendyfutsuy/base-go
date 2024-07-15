package usecase

import (
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category"
	cobsubcob "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob"
)

type cobsubcobUsecase struct {
	categoryRepo   category.Repository
	cobsubcobRepo     cobsubcob.Repository
	contextTimeout time.Duration
}

func NewCobSubcobUsecase(cr category.Repository, sr cobsubcob.Repository, timeout time.Duration) cobsubcob.Usecase {
	return &cobsubcobUsecase{
		categoryRepo:   cr,
		cobsubcobRepo:     sr,
		contextTimeout: timeout,
	}
}
