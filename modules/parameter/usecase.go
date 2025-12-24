package parameter

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateParameter, authId string) (*models.Parameter, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateParameter, authId string) (*models.Parameter, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Parameter, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqParameterIndexFilter) ([]models.Parameter, int, error)
	GetAll(ctx context.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error)
	Export(ctx context.Context, filter dto.ReqParameterIndexFilter) ([]byte, error)
}
