package type_module

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/type/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateType, authId string) (*models.Type, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateType, authId string) (*models.Type, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Type, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqTypeIndexFilter) ([]models.Type, int, error)
	GetAll(ctx context.Context, filter dto.ReqTypeIndexFilter) ([]models.Type, error)
	Export(ctx context.Context, filter dto.ReqTypeIndexFilter) ([]byte, error)
}
