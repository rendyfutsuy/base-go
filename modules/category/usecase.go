package category

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"github.com/labstack/echo/v4"
)

type Usecase interface {
	CreateCategory(c echo.Context, req *dto.ReqCreateCategory, authId string) (categoryRes *models.Category, err error)
	GetCategoryByID(id string) (category *models.Category, err error)
	GetAllCategory() (categoryes []models.Category, err error)
	GetIndexCategory(req request.PageRequest) (categoryes []models.Category, total int, err error)
	UpdateCategory(id string, req *dto.ReqUpdateCategory, authId string) (categoryRes *models.Category, err error)
	SoftDeleteCategory(id string, authId string) (categoryRes *models.Category, err error)
}
