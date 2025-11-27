package sub_group

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/sub-group/dto"
)

type Usecase interface {
	Create(c echo.Context, req *dto.ReqCreateSubGroup, authId string) (*models.SubGroup, error)
	Update(c echo.Context, id string, req *dto.ReqUpdateSubGroup, authId string) (*models.SubGroup, error)
	Delete(c echo.Context, id string, authId string) error
	GetByID(c echo.Context, id string) (*models.SubGroup, error)
	GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error)
	GetAll(c echo.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error)
	Export(c echo.Context, filter dto.ReqSubGroupIndexFilter) ([]byte, error)
	ExistsInTypes(c echo.Context, subGroupID string) (bool, error)
}
