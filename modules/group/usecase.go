package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
)

type Usecase interface {
	Create(c echo.Context, req *dto.ReqCreateGroup, authId string) (*models.GoodsGroup, error)
	Update(c echo.Context, id string, req *dto.ReqUpdateGroup, authId string) (*models.GoodsGroup, error)
	Delete(c echo.Context, id string, authId string) error
	GetByID(c echo.Context, id string) (*models.GoodsGroup, error)
	GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error)
	GetAll(c echo.Context, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, error)
	Export(c echo.Context, filter dto.ReqGroupIndexFilter) ([]byte, error)
	ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error)
}
