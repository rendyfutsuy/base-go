package type_module

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/type/dto"
)

type Usecase interface {
	Create(c echo.Context, req *dto.ReqCreateType, authId string) (*models.Type, error)
	Update(c echo.Context, id string, req *dto.ReqUpdateType, authId string) (*models.Type, error)
	Delete(c echo.Context, id string, authId string) error
	GetByID(c echo.Context, id string) (*models.Type, error)
	GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqTypeIndexFilter) ([]models.Type, int, error)
	GetAll(c echo.Context, filter dto.ReqTypeIndexFilter) ([]models.Type, error)
	Export(c echo.Context, filter dto.ReqTypeIndexFilter) ([]byte, error)
}
