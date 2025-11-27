package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
	"gorm.io/gorm"
)

type backingRepository struct {
	DB *gorm.DB
}

func NewBackingRepository(db *gorm.DB) *backingRepository {
	return &backingRepository{DB: db}
}

func (r *backingRepository) Create(ctx context.Context, typeID uuid.UUID, name string, createdBy string) (*models.Backing, error) {
	now := time.Now().UTC()
	b := &models.Backing{
		TypeID:    typeID,
		Name:      name,
		CreatedAt: now,
		CreatedBy: createdBy,
		UpdatedAt: now,
		UpdatedBy: createdBy,
	}
	// Omit backing_code to let database generate it using DEFAULT generate_backing_code()
	if err := r.DB.WithContext(ctx).Omit("backing_code").Create(b).Error; err != nil {
		return nil, err
	}
	// if b not update, return error
	if b.ID == uuid.Nil {
		return nil, errors.New(constants.BackingCreateFailedIDNotSet)
	}
	return b, nil
}

func (r *backingRepository) Update(ctx context.Context, id uuid.UUID, typeID uuid.UUID, name string, updatedBy string) (*models.Backing, error) {
	updates := map[string]interface{}{
		"type_id":    typeID,
		"name":       name,
		"updated_at": time.Now().UTC(),
		"updated_by": updatedBy,
	}
	b := &models.Backing{}
	err := r.DB.WithContext(ctx).Model(&models.Backing{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return b, nil
}

func (r *backingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Backing{}).Error
}

func (r *backingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Backing, error) {
	b := &models.Backing{}
	query := r.DB.WithContext(ctx).
		Table("backings b").
		Select(`
			b.id,
			b.type_id,
			b.backing_code,
			b.name,
			b.created_at,
			b.created_by,
			b.updated_at,
			b.updated_by,
			t.name as type_name,
			t.subgroup_id as subgroup_id,
			sg.name as subgroup_name,
			sg.goods_group_id as goods_group_id,
			gg.name as group_name
		`).
		Joins("LEFT JOIN types t ON b.type_id = t.id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN sub_groups sg ON t.subgroup_id = sg.id AND sg.deleted_at IS NULL").
		Joins("LEFT JOIN goods_group gg ON sg.goods_group_id = gg.id AND gg.deleted_at IS NULL").
		Where("b.id = ? AND b.deleted_at IS NULL", id)

	err := query.Scan(b).Error
	if err != nil {
		return nil, err
	}
	// Scan() doesn't return error for record not found, so check if ID is nil
	if b.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}
	return b, nil
}

func (r *backingRepository) ExistsByNameInType(ctx context.Context, typeID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Backing{}).Where("type_id = ? AND name = ?", typeID, name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *backingRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error) {
	var backings []models.Backing
	query := r.DB.WithContext(ctx).
		Table("backings b").
		Select(`
			b.id,
			b.type_id,
			b.backing_code,
			b.name,
			b.created_at,
			b.updated_at,
			t.name as type_name,
			sg.name as subgroup_name,
			gg.name as group_name
		`).
		Joins("LEFT JOIN types t ON b.type_id = t.id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN sub_groups sg ON t.subgroup_id = sg.id AND sg.deleted_at IS NULL").
		Joins("LEFT JOIN goods_group gg ON sg.goods_group_id = gg.id AND gg.deleted_at IS NULL").
		Where("b.deleted_at IS NULL")

	// Apply search from PageRequest
	searchQuery := req.Search
	query = request.ApplySearchCondition(query, searchQuery, []string{"b.name", "b.backing_code"})

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "b.created_at",
		DefaultSortOrder:   "DESC",
		MaxPerPage:         100,
		SortMapping:        mapBackingIndexSortColumn,
		NaturalSortColumns: []string{"b.name"}, // Enable natural sorting for b.name
	}, &backings)
	if err != nil {
		return nil, 0, err
	}
	return backings, total, nil
}

func (r *backingRepository) GetAll(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error) {
	var backings []models.Backing
	query := r.DB.WithContext(ctx).
		Table("backings b").
		Select(`
			b.id,
			b.type_id,
			b.backing_code,
			b.name,
			b.created_at,
			b.updated_at,
			t.name as type_name,
			sg.name as subgroup_name,
			gg.name as group_name
		`).
		Joins("LEFT JOIN types t ON b.type_id = t.id AND t.deleted_at IS NULL").
		Joins("LEFT JOIN sub_groups sg ON t.subgroup_id = sg.id AND sg.deleted_at IS NULL").
		Joins("LEFT JOIN goods_group gg ON sg.goods_group_id = gg.id AND gg.deleted_at IS NULL").
		Where("b.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchCondition(query, filter.Search, []string{"b.name", "b.backing_code"})

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Determine sorting
	sortBy := "b.created_at"
	if mapped := mapBackingIndexSortColumn(filter.SortBy); mapped != "" {
		sortBy = mapped
	}

	// Determine sorting with natural sorting support
	sortExpression := request.BuildSortExpressionForExport(
		sortBy,
		filter.SortOrder,
		"b.created_at",
		"DESC",
		mapBackingIndexSortColumn,
		[]string{"b.name"}, // Enable natural sorting for backing name
	)

	// Order results
	// Use Scan() for JOIN queries with custom SELECT
	if err := query.Order(sortExpression).Scan(&backings).Error; err != nil {
		return nil, err
	}
	return backings, nil
}
