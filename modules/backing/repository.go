package backing

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
)

type Repository interface {
	Create(ctx context.Context, typeID uuid.UUID, name string, createdBy string) (*models.Backing, error)
	Update(ctx context.Context, id uuid.UUID, typeID uuid.UUID, name string, updatedBy string) (*models.Backing, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Backing, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error)
	GetAll(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error)
	ExistsByNameInType(ctx context.Context, typeID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)
}
