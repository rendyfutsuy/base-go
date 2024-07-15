package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (u *categoryUsecase) CreateCategory(c echo.Context, req *dto.ReqCreateCategory, authId string) (categoryRes *models.Category, err error) {

	count, err := u.categoryRepo.CountCategory()
	if err != nil {
		return nil, err
	}

	formatCount := fmt.Sprintf("%07d", *count+1)

	categoryDb := req.ToDBCreateCategory(formatCount, authId)

	categoryRes, err = u.categoryRepo.CreateCategory(nil, categoryDb)
	if err != nil {
		return nil, err
	}

	return categoryRes, err
}

func (u *categoryUsecase) GetCategoryByID(id string) (category *models.Category, err error) {

	uId := uuid.MustParse(id)

	return u.categoryRepo.GetCategoryByID(uId)
}

func (u *categoryUsecase) GetIndexCategory(req request.PageRequest) (categoryes []models.Category, total int, err error) {
	return u.categoryRepo.GetIndexCategory(req)
}

func (u *categoryUsecase) GetAllCategory() (categoryes []models.Category, err error) {

	return u.categoryRepo.GetAllCategory()
}

func (u *categoryUsecase) UpdateCategory(id string, req *dto.ReqUpdateCategory, authId string) (categoryRes *models.Category, err error) {

	uId := uuid.MustParse(id)

	categoryDb := req.ToDBUpdateCategory(authId)

	return u.categoryRepo.UpdateCategory(uId, categoryDb)
}

func (u *categoryUsecase) SoftDeleteCategory(id string, authId string) (categoryRes *models.Category, err error) {

	uId := uuid.MustParse(id)

	categoryDb := dto.ToDBDeleteCategory{
		DeletedByID: authId,
	}

	return u.categoryRepo.SoftDeleteCategory(uId, categoryDb)
}
