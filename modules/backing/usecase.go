package backing

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
)

type Usecase interface {
	Create(c echo.Context, req *dto.ReqCreateBacking, authId string) (*models.Backing, error)
	Update(c echo.Context, id string, req *dto.ReqUpdateBacking, authId string) (*models.Backing, error)
	Delete(c echo.Context, id string, authId string) error
	GetByID(c echo.Context, id string) (*models.Backing, error)
	GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error)
	GetAll(c echo.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error)
	Export(c echo.Context, filter dto.ReqBackingIndexFilter) ([]byte, error)
}
