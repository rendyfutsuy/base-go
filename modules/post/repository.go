package post

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/post/dto"
)

type Repository interface {
	Create(ctx context.Context, createdBy uuid.UUID, data dto.ToDBPost) (*models.Post, error)
	Update(ctx context.Context, id uuid.UUID, data dto.ToDBPost) (*models.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqPostIndexFilter) ([]models.Post, int, error)
	GetAll(ctx context.Context, filter dto.ReqPostIndexFilter) ([]models.Post, error)
}
