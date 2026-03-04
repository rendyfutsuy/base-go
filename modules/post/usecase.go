package post

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/post/dto"
)

type Usecase interface {
	Create(ctx context.Context, req *dto.ReqCreatePost, authId string, thumbnailData []byte, thumbnailName string) (*models.Post, error)
	Update(ctx context.Context, id string, req *dto.ReqUpdatePost, authId string, thumbnailData []byte, thumbnailName string) (*models.Post, error)
	Delete(ctx context.Context, id string, authId string) error
	GetByID(ctx context.Context, id string) (*models.Post, error)
	GetParameterReferences(ctx context.Context, id string) (*dto.ReferenceObject, *dto.ReferenceObject, []dto.ReferenceObject, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqPostIndexFilter) ([]models.Post, int, error)
	GetAll(ctx context.Context, filter dto.ReqPostIndexFilter) ([]models.Post, error)
}
