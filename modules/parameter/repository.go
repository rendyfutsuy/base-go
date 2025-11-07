package parameter

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
)

type Repository interface {
	Create(ctx context.Context, code, name string, value, typeVal, desc *string) (*models.Parameter, error)
	Update(ctx context.Context, id uuid.UUID, code, name string, value, typeVal, desc *string) (*models.Parameter, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Parameter, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqParameterIndexFilter) ([]models.Parameter, int, error)
	GetAll(ctx context.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error)
	ExistsByCode(ctx context.Context, code string, excludeID uuid.UUID) (bool, error)
	ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
}
