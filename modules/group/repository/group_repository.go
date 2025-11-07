package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
	"gorm.io/gorm"
)

type groupRepository struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{DB: db}
}

func (r *groupRepository) Create(ctx context.Context, name string) (*models.GoodsGroup, error) {
	now := time.Now().UTC()
	gg := &models.GoodsGroup{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	// Omit group_code to let database generate it using DEFAULT generate_group_code()
	if err := r.DB.WithContext(ctx).Omit("group_code").Create(gg).Error; err != nil {
		return nil, err
	}
	// if gg not update, return error
	if gg.ID == uuid.Nil {
		return nil, errors.New(constants.GroupCreateFailedIDNotSet)
	}
	return gg, nil
}

func (r *groupRepository) Update(ctx context.Context, id uuid.UUID, name string) (*models.GoodsGroup, error) {
	updates := map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().UTC(),
	}
	gg := &models.GoodsGroup{}
	err := r.DB.WithContext(ctx).Model(&models.GoodsGroup{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(gg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return gg, nil
}

func (r *groupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.GoodsGroup{}).Error
}

func (r *groupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.GoodsGroup, error) {
	gg := &models.GoodsGroup{}
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(gg).Error; err != nil {
		return nil, err
	}
	return gg, nil
}

func (r *groupRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.GoodsGroup{}).Where("name = ?", name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error) {
	var groups []models.GoodsGroup
	query := r.DB.WithContext(ctx).Table("goods_group gg").Select("gg.id, gg.group_code, gg.name, gg.created_at, gg.updated_at").
		Where("gg.deleted_at IS NULL")

	// Apply search from PageRequest
	searchQuery := req.Search
	query = request.ApplySearchCondition(query, searchQuery, []string{"gg.name", "gg.group_code"})

	// Apply filter conditions (can be extended in the future)
	// Example: if len(filter.GroupCodes) > 0 { query = query.Where("gg.group_code IN (?)", filter.GroupCodes) }

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:    "gg.created_at",
		DefaultSortOrder: "DESC",
		AllowedColumns:   []string{"id", "group_code", "name", "created_at", "updated_at"},
		ColumnPrefix:     "gg.",
		MaxPerPage:       100,
	}, &groups)
	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (r *groupRepository) GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.GoodsGroup, error) {
	var groups []models.GoodsGroup
	query := r.DB.WithContext(ctx).Table("goods_group gg").Select("gg.id, gg.group_code, gg.name, gg.created_at, gg.updated_at").
		Where("gg.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchCondition(query, filter.Search, []string{"gg.name", "gg.group_code"})

	// Apply filter conditions (can be extended in the future)
	// Example: if len(filter.GroupCodes) > 0 { query = query.Where("gg.group_code IN (?)", filter.GroupCodes) }

	// Order by created_at DESC (no pagination)
	if err := query.Order("gg.created_at DESC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
