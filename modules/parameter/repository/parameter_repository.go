package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
	rsearchparam "github.com/rendyfutsuy/base-go/modules/parameter/repository/searches"
	"gorm.io/gorm"
)

type parameterRepository struct {
	DB *gorm.DB
}

func NewParameterRepository(db *gorm.DB) *parameterRepository {
	return &parameterRepository{
		DB: db,
	}
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
	query := r.DB.WithContext(ctx).Table("parameters p").Select("p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

		// Apply search from PageRequest
	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, rsearchparam.NewParameterSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "p.created_at",
		DefaultSortOrder:   "DESC",
		MaxPerPage:         100,
		SortMapping:        mapParameterIndexSortColumn,
		NaturalSortColumns: []string{"p.name"}, // Enable natural sorting for p.name
	}, &parameters)
	if err != nil {
		return nil, 0, err
	}
	return parameters, total, nil
}

func (r *parameterRepository) GetAll(ctx context.Context, filter dto.ReqParameterIndexFilter) ([]models.Parameter, error) {
	var parameters []models.Parameter
	query := r.DB.WithContext(ctx).Table("parameters p").Select("p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchConditionFromInterface(query, filter.Search, rsearchparam.NewParameterSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Determine sorting
	sortBy := "p.created_at"
	if mapped := mapParameterIndexSortColumn(filter.SortBy); mapped != "" {
		sortBy = mapped
	}

	sortOrder := request.ValidateAndSanitizeSortOrder(filter.SortOrder)
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	// Order results
	if err := query.Order(sortBy + " " + sortOrder).Find(&parameters).Error; err != nil {
		return nil, err
	}
	return parameters, nil
}
