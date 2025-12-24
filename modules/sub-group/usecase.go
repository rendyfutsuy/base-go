package sub_group

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/sub-group/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateSubGroup, authId string) (*models.SubGroup, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateSubGroup, authId string) (*models.SubGroup, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.SubGroup, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error)
	GetAll(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error)
	Export(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]byte, error)
	ExistsInTypes(ctx context.Context, subGroupID string) (bool, error)
}
