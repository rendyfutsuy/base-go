package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

// GetPermissionGroupByID retrieves an permission_group information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an PermissionGroup pointer and an error.
func (repo *roleRepository) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	permissionGroup = &models.PermissionGroup{}

	// Use Raw query with parameter binding for ARRAY_AGG
	err = repo.DB.WithContext(ctx).
		Raw(`
			SELECT
				pg.id AS permission_group_id,
				pg.name AS permission_group_name,
				ARRAY_AGG(p.name) AS permissions,
				pg.created_at,
				pg.updated_at,
				pg.deleted_at
			FROM
				permission_groups pg
			LEFT JOIN
				permissions_modules ppg
			ON
				pg.id = ppg.permission_group_id
			LEFT JOIN
				permissions p
			ON
				ppg.permission_id = p.id
			WHERE
				pg.id = ? AND pg.deleted_at IS NULL
			GROUP BY
				pg.id, pg.name
		`, id).
		Scan(permissionGroup).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission_group permission_group with id %s not found", id)
		}
		return nil, err
	}

	// Handle NULL arrays from LEFT JOIN
	if permissionGroup.PermissionNames == nil {
		permissionGroup.PermissionNames = []utils.NullString{}
	}

	return permissionGroup, nil
}

// GetIndexPermissionGroup retrieves a paginated list of permission_group information from the database.
func (repo *roleRepository) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (permissionGroups []models.PermissionGroup, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Build base query
	query := repo.DB.WithContext(ctx).
		Table("permission_groups permission_group").
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("permission_group.deleted_at IS NULL")

	// Apply search with parameter binding
	if searchQuery != "" {
		query = query.Where("permission_group.name ILIKE ?", "%"+searchQuery+"%")
	}

	// Count total (before pagination)
	countQuery := query
	var totalCount int64
	err = countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	total = int(totalCount)

	// Apply sorting
	sortBy := "permission_group.created_at"
	sortOrder := "DESC"
	if req.SortBy != "" {
		sortBy = req.SortBy
		if req.SortOrder != "" {
			sortOrder = req.SortOrder
		}
	}

	// Apply pagination and sorting
	err = query.
		Order(sortBy + " " + sortOrder).
		Limit(req.PerPage).
		Offset(offSet).
		Find(&permissionGroups).Error

	if err != nil {
		return nil, 0, err
	}

	return permissionGroups, total, nil
}

// GetAllPermissionGroup retrieves all permission_group information entries from the database.
//
// Returns a slice of models.PermissionGroup and an error.
func (repo *roleRepository) GetAllPermissionGroup(ctx context.Context) (permissionGroups []models.PermissionGroup, err error) {
	err = repo.DB.WithContext(ctx).
		Table("permission_groups permission_group").
		Select("id", "name", "module", "created_at", "updated_at", "deleted_at").
		Where("permission_group.deleted_at IS NULL").
		Order("permission_group.module ASC").
		Find(&permissionGroups).Error

	if err != nil {
		return nil, err
	}

	return permissionGroups, nil
}

// CountPermissionGroup retrieves the count of permission_group information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountPermissionGroup(ctx context.Context) (count *int, err error) {
	var result int64
	err = repo.DB.WithContext(ctx).
		Model(&models.PermissionGroup{}).
		Count(&result).Error

	if err != nil {
		return nil, err
	}

	resultInt := int(result)
	count = &resultInt
	return count, nil
}

// PermissionGroupNameIsNotDuplicated checks if the provided permission_group name is not duplicated in the database.
func (repo *roleRepository) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.PermissionGroup{}).
		Where("name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetDuplicatedPermissionGroup retrieves the permission_group information with the given name and excluded ID from the database.
func (repo *roleRepository) GetDuplicatedPermissionGroup(ctx context.Context, name string, excludedId uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	permissionGroup = &models.PermissionGroup{}

	query := repo.DB.WithContext(ctx).
		Select("id", "name", "created_at", "updated_at").
		Where("name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(permissionGroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return permissionGroup, nil
}
