package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
	"gorm.io/gorm"
)

type parameterRepository struct {
	DB *gorm.DB
}

func NewParameterRepository(db *gorm.DB) *parameterRepository {
	return &parameterRepository{DB: db}
}

func (r *parameterRepository) Create(ctx context.Context, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	now := time.Now().UTC()
	p := &models.Parameter{
		Code:        code,
		Name:        name,
		Value:       value,
		Type:        typeVal,
		Description: desc,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := r.DB.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	if p.ID == uuid.Nil {
		return nil, errors.New("failed to create parameter: ID not set")
	}
	return p, nil
}

func (r *parameterRepository) Update(ctx context.Context, id uuid.UUID, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	updates := map[string]interface{}{
		"code":       code,
		"name":       name,
		"updated_at": time.Now().UTC(),
	}
	if value != nil {
		updates["value"] = *value
	} else {
		updates["value"] = nil
	}
	if typeVal != nil {
		updates["type"] = *typeVal
	} else {
		updates["type"] = nil
	}
	if desc != nil {
		updates["description"] = *desc
	} else {
		updates["description"] = nil
	}
	p := &models.Parameter{}
	err := r.DB.WithContext(ctx).Model(&models.Parameter{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return p, nil
}

func (r *parameterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Parameter{}).Error
}

func (r *parameterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Parameter, error) {
	p := &models.Parameter{}
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *parameterRepository) ExistsByCode(ctx context.Context, code string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Parameter{}).Where("code = ?", code)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *parameterRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Parameter{}).Where("name = ?", name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *parameterRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqParameterIndexFilter) ([]models.Parameter, int, error) {
	var parameters []models.Parameter
	query := r.DB.WithContext(ctx).Table("parameter p").Select("p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

	// Apply search from PageRequest
	searchQuery := req.Search
	query = request.ApplySearchCondition(query, searchQuery, []string{"p.name", "p.code"})

	// Apply filter conditions
	if len(filter.Types) > 0 {
		query = query.Where("p.type IN ?", filter.Types)
	}
	if len(filter.Names) > 0 {
		query = query.Where("p.name IN ?", filter.Names)
	}

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:    "p.created_at",
		DefaultSortOrder: "DESC",
		AllowedColumns:   []string{"id", "code", "name", "value", "type", "created_at", "updated_at"},
		ColumnPrefix:     "p.",
		MaxPerPage:       100,
	}, &parameters)
	if err != nil {
		return nil, 0, err
	}
	return parameters, total, nil
}

func (r *parameterRepository) GetAll(ctx context.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error) {
	var parameters []models.Parameter
	query := r.DB.WithContext(ctx).Table("parameter p").Select("p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchCondition(query, filter.Search, []string{"p.name", "p.code"})

	// Apply filter conditions
	if len(filter.Types) > 0 {
		query = query.Where("p.type IN ?", filter.Types)
	}
	if len(filter.Names) > 0 {
		query = query.Where("p.name IN ?", filter.Names)
	}

	// Order by created_at DESC (no pagination)
	if err := query.Order("p.created_at DESC").Find(&parameters).Error; err != nil {
		return nil, err
	}
	return parameters, nil
}
