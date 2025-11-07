package parameter

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
)

type Usecase interface {
	Create(c echo.Context, req *dto.ReqCreateParameter, authId string) (*models.Parameter, error)
	Update(c echo.Context, id string, req *dto.ReqUpdateParameter, authId string) (*models.Parameter, error)
	Delete(c echo.Context, id string, authId string) error
	GetByID(c echo.Context, id string) (*models.Parameter, error)
	GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqParameterIndexFilter) ([]models.Parameter, int, error)
	GetAll(c echo.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error)
	Export(c echo.Context, filter dto.ReqParameterIndexFilter) ([]byte, error)
}
