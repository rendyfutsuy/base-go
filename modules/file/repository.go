package file

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/file/dto"
)

type Repository interface {
	Create(ctx context.Context, data dto.ToDBFile) (*models.File, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AssignFilesToModule(ctx context.Context, moduleType string, moduleID uuid.UUID, fileIDs []uuid.UUID) error
	RemoveFilesFromModule(ctx context.Context, moduleType string, moduleID uuid.UUID) error
}
