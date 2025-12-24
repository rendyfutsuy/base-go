package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreateGroup, authId string) (*models.GoodsGroup, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdateGroup, authId string) (*models.GoodsGroup, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.GoodsGroup, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error)
	GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, error)
	Export(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]byte, error)
	ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error)
}
