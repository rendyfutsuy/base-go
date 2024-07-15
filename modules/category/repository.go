package category

import (
	"database/sql"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// crud
	CreateCategory(trx *sql.Tx, categoryReq dto.ToDBCreateCategory) (categoryRes *models.Category, err error)
	GetCategoryByID(id uuid.UUID) (category *models.Category, err error)
	GetIndexCategory(req request.PageRequest) (categoryes []models.Category, total int, err error)
	GetAllCategory() (categoryes []models.Category, err error)
	UpdateCategory(id uuid.UUID, categoryReq dto.ToDBUpdateCategory) (categoryRes *models.Category, err error)
	SoftDeleteCategory(id uuid.UUID, categoryReq dto.ToDBDeleteCategory) (categoryRes *models.Category, err error)

	// general
	CountCategory() (count *int, err error)
}
