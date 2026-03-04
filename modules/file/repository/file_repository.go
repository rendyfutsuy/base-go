package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/file/dto"
	"gorm.io/gorm"
)

type fileRepository struct {
	DB *gorm.DB
}

func NewFileRepository(db *gorm.DB) *fileRepository {
	return &fileRepository{DB: db}
}

func (r *fileRepository) Create(ctx context.Context, data dto.ToDBFile) (*models.File, error) {
	now := time.Now().UTC()
	f := &models.File{
		Name:        data.Name,
		FilePath:    data.FilePath,
		Description: data.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := r.DB.WithContext(ctx).Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.File{}).Error
}

func (r *fileRepository) AssignFilesToModule(ctx context.Context, input dto.AssignFilesToModule) error {
	if len(input.Items) == 0 {
		return nil
	}
	now := time.Now().UTC()
	items := make([]models.FilesToModule, 0, len(input.Items))
	for _, it := range input.Items {
		tp := ""
		if it.Type != nil {
			tp = *it.Type
		}
		items = append(items, models.FilesToModule{
			FileID:     it.FileID,
			ModuleType: input.ModuleType,
			ModuleID:   input.ModuleID,
			Type:       tp,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}
	return r.DB.WithContext(ctx).Create(&items).Error
}

func (r *fileRepository) RemoveFilesFromModule(ctx context.Context, moduleType string, moduleID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Where("module_type = ? AND module_id = ?", moduleType, moduleID).
		Delete(&models.FilesToModule{}).Error
}
