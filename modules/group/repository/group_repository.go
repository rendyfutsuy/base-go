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
	rsearchgroup "github.com/rendyfutsuy/base-go/modules/group/repository/searches"
	"gorm.io/gorm"
)

type groupRepository struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{
		DB: db,
	}
}

func (r *groupRepository) Create(ctx context.Context, name string, createdBy string) (*models.Group, error) {
	now := time.Now().UTC()
	gg := &models.Group{
		Name:      name,
		CreatedAt: now,
		CreatedBy: createdBy,
		UpdatedAt: now,
		UpdatedBy: createdBy,
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

func (r *groupRepository) Update(ctx context.Context, id uuid.UUID, name string, updatedBy string) (*models.Group, error) {
	updates := map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().UTC(),
		"updated_by": updatedBy,
	}
	err := r.DB.WithContext(ctx).Model(&models.Group{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	// Get updated group with deletable status
	gg := &models.Group{}
	err = r.DB.WithContext(ctx).Table("groups gg").
		Select(`
			gg.id, 
			gg.group_code, 
			gg.name, 
			gg.created_at, 
			gg.updated_at,
			NOT EXISTS (
				SELECT 1 
				FROM sub_groups sg 
				WHERE sg.groups_id = gg.id 
				AND sg.deleted_at IS NULL
			) as deletable
		`).
		Where("gg.id = ? AND gg.deleted_at IS NULL", id).
		First(gg).Error
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func (r *groupRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now().UTC(),
		"deleted_by": deletedBy,
	}
	return r.DB.WithContext(ctx).Model(&models.Group{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *groupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	gg := &models.Group{}
	err := r.DB.WithContext(ctx).Table("groups gg").
		Select(`
			gg.id, 
			gg.group_code, 
			gg.name, 
			gg.created_at, 
			gg.updated_at,
			NOT EXISTS (
				SELECT 1 
				FROM sub_groups sg 
				WHERE sg.groups_id = gg.id 
				AND sg.deleted_at IS NULL
			) as deletable
		`).
		Where("gg.id = ? AND gg.deleted_at IS NULL", id).
		First(gg).Error
	if err != nil {
		return nil, err
	}
	return gg, nil
}

func (r *groupRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Group{}).Where("name = ?", name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *groupRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.Group, int, error) {
	var groups []models.Group
	query := r.DB.WithContext(ctx).Table("groups gg").
		Select(`
			gg.id, 
			gg.group_code, 
			gg.name, 
			gg.created_at, 
			gg.updated_at,
			NOT EXISTS (
				SELECT 1 
				FROM sub_groups sg 
				WHERE sg.groups_id = gg.id 
				AND sg.deleted_at IS NULL
			) as deletable
		`).
		Where("gg.deleted_at IS NULL")

		// Apply search from PageRequest
	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, rsearchgroup.NewGroupSearchHelper())

	// Apply filter conditions (can be extended in the future)
	// Example: if len(filter.GroupCodes) > 0 { query = query.Where("gg.group_code IN (?)", filter.GroupCodes) }

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "gg.created_at",
		DefaultSortOrder:   "DESC",
		MaxPerPage:         100,
		SortMapping:        mapGroupIndexSortColumn,
		NaturalSortColumns: []string{"gg.name"}, // Enable natural sorting for gg.name
	}, &groups)
	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (r *groupRepository) GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.Group, error) {
	var groups []models.Group
	query := r.DB.WithContext(ctx).Table("groups gg").Select("gg.id, gg.group_code, gg.name, gg.created_at, gg.updated_at").
		Where("gg.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchConditionFromInterface(query, filter.Search, rsearchgroup.NewGroupSearchHelper())

	// Apply filter conditions (can be extended in the future)
	// Example: if len(filter.GroupCodes) > 0 { query = query.Where("gg.group_code IN (?)", filter.GroupCodes) }

	// Determine sorting with natural sorting support
	sortExpression := request.BuildSortExpressionForExport(
		filter.SortBy,
		filter.SortOrder,
		"gg.created_at",
		"DESC",
		mapGroupIndexSortColumn,
		[]string{"gg.name"}, // Enable natural sorting for group name
	)

	// Order results
	if err := query.Order(sortExpression).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *groupRepository) ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.SubGroup{}).
		Where("groups_id = ? AND deleted_at IS NULL", groupID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
