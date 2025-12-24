package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
)

type Repository interface {
	Create(ctx context.Context, name string, createdBy string) (*models.Group, error)
	Update(ctx context.Context, id uuid.UUID, name string, updatedBy string) (*models.Group, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.Group, int, error)
	GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.Group, error)
	ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
	ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error)
}
