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

func (r *parameterRepository) AssignParametersToModule(ctx context.Context, moduleType string, moduleID uuid.UUID, parameterIDs []uuid.UUID) error {
	if len(parameterIDs) == 0 {
		return nil
	}
	now := time.Now().UTC()
	items := make([]models.ParametersToModule, 0, len(parameterIDs))
	for _, pid := range parameterIDs {
		items = append(items, models.ParametersToModule{
			ParameterID: pid,
			ModuleType:  moduleType,
			ModuleID:    moduleID,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return r.DB.WithContext(ctx).Create(&items).Error
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

func (r *parameterRepository) SetParent(ctx context.Context, id uuid.UUID, parentID uuid.UUID) error {
	updates := map[string]interface{}{
		"parent_id":  parentID,
		"updated_at": time.Now().UTC(),
	}
	return r.DB.WithContext(ctx).Model(&models.Parameter{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *parameterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Parameter{}).Error
}

func (r *parameterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Parameter, error) {
	p := &models.Parameter{}
	if err := r.DB.WithContext(ctx).
		Table("parameters p").
		Select(`p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at, p.parent_id, parent.name AS parent_name,
			(
				NOT EXISTS (SELECT 1 FROM parameters_to_module ptm WHERE ptm.parameter_id = p.id)
				AND NOT EXISTS (SELECT 1 FROM parameters child WHERE child.parent_id = p.id AND child.deleted_at IS NULL)
			) AS deletable`).
		Joins("LEFT JOIN parameters parent ON parent.id = p.parent_id AND parent.deleted_at IS NULL").
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		Scan(p).Error; err != nil {
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
	query := r.DB.WithContext(ctx).Table("parameters p").Select(`p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at,
		(
			NOT EXISTS (SELECT 1 FROM parameters_to_module ptm WHERE ptm.parameter_id = p.id)
			AND NOT EXISTS (SELECT 1 FROM parameters child WHERE child.parent_id = p.id AND child.deleted_at IS NULL)
		) AS deletable`).
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
	query := r.DB.WithContext(ctx).Table("parameters p").Select(`p.id, p.code, p.name, p.value, p.type, p.description, p.created_at, p.updated_at,
		(
			NOT EXISTS (SELECT 1 FROM parameters_to_module ptm WHERE ptm.parameter_id = p.id)
			AND NOT EXISTS (SELECT 1 FROM parameters child WHERE child.parent_id = p.id AND child.deleted_at IS NULL)
		) AS deletable`).
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

func (r *parameterRepository) GetByModule(ctx context.Context, moduleType string, moduleID uuid.UUID) ([]models.Parameter, error) {
	var params []models.Parameter
	if err := r.DB.WithContext(ctx).
		Table("parameters p").
		Select("p.id, p.name, p.type").
		Joins("JOIN parameters_to_module ptm ON ptm.parameter_id = p.id").
		Where("ptm.module_type = ? AND ptm.module_id = ? AND p.deleted_at IS NULL", moduleType, moduleID).
		Find(&params).Error; err != nil {
		return nil, err
	}
	return params, nil
}
