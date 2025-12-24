package backing

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateBacking, authId string) (*models.Backing, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateBacking, authId string) (*models.Backing, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Backing, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error)
	GetAll(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error)
	Export(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]byte, error)
}
