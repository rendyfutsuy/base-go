package expedition

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/expedition/dto"
)

// CreateExpeditionParams contains parameters for creating an expedition
type CreateExpeditionParams struct {
	ExpeditionName string
	Address        string
	TelpNumbers    []dto.TelpNumberItem
	PhoneNumbers   []string
	Notes          *string
	CreatedBy      string
}

// UpdateExpeditionParams contains parameters for updating an expedition
type UpdateExpeditionParams struct {
	ExpeditionName string
	Address        string
	TelpNumbers    []dto.TelpNumberItem
	PhoneNumbers   []string
	Notes          *string
	UpdatedBy      string
}

type Repository interface {
	Create(ctx context.Context, params CreateExpeditionParams) (*models.Expedition, error)
	Update(ctx context.Context, id uuid.UUID, params UpdateExpeditionParams) (*models.Expedition, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Expedition, error)
	GetContactsByExpeditionID(ctx context.Context, expeditionID uuid.UUID) ([]models.ExpeditionContact, error)
	GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error)
	GetAll(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, error)
	GetAllForExport(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]dto.ExpeditionExport, error)
	ExistsByExpeditionName(ctx context.Context, expeditionName string, excludeID uuid.UUID) (bool, error)
}
