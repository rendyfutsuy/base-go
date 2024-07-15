package class

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/dto"
	"github.com/google/uuid"
)

type Repository interface {
	// migration
	CreateTable(sqlFilePath string) (err error)

	// crud
	CreateClass(classReq dto.ToDBCreateClass) (classRes *models.Class, err error)
	GetClassByID(id uuid.UUID) (class *models.Class, err error)
	GetIndexClass(req request.PageRequest) (classes []models.Class, total int, err error)
	GetAllClass() (classes []models.Class, err error)
	UpdateClass(id uuid.UUID, classReq dto.ToDBUpdateClass) (classRes *models.Class, err error)
	SoftDeleteClass(id uuid.UUID, classReq dto.ToDBDeleteClass) (classRes *models.Class, err error)

	// general
	CountClass() (count *int, err error)
}
