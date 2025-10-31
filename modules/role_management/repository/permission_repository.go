package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"gorm.io/gorm"
)

// GetPermissionByID retrieves an permission information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an Permission pointer and an error.
func (repo *roleRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (permission *models.Permission, err error) {
	permission = &models.Permission{}

	err = repo.DB.WithContext(ctx).
		Table("permissions permission").
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("permission.id = ? AND permission.deleted_at IS NULL", id).
		First(permission).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission permission with id %s not found", id)
		}
		return nil, err
	}

	return permission, nil
}

// GetIndexPermission retrieves a paginated list of permission information from the database.
func (repo *roleRepository) GetIndexPermission(ctx context.Context, req request.PageRequest) (permissions []models.Permission, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Build base query
	query := repo.DB.WithContext(ctx).
		Table("permissions permission").
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("permission.deleted_at IS NULL")

	// Apply search with parameter binding
	if searchQuery != "" {
		query = query.Where("permission.name ILIKE ?", "%"+searchQuery+"%")
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
	sortBy := "permission.created_at"
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
		Find(&permissions).Error

	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetAllPermission retrieves all permission information entries from the database.
//
// Returns a slice of models.Permission and an error.
func (repo *roleRepository) GetAllPermission(ctx context.Context) (permissions []models.Permission, err error) {
	err = repo.DB.WithContext(ctx).
		Table("permissions permission").
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("permission.deleted_at IS NULL").
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// CountPermission retrieves the count of permission information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountPermission(ctx context.Context) (count *int, err error) {
	var result int64
	err = repo.DB.WithContext(ctx).
		Model(&models.Permission{}).
		Count(&result).Error

	if err != nil {
		return nil, err
	}

	resultInt := int(result)
	count = &resultInt
	return count, nil
}

// PermissionNameIsNotDuplicated checks if the provided permission name is not duplicated in the database.
func (repo *roleRepository) PermissionNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.Permission{}).
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

// GetDuplicatedPermission retrieves the permission information with the given name and excluded ID from the database.
func (repo *roleRepository) GetDuplicatedPermission(ctx context.Context, name string, excludedId uuid.UUID) (permission *models.Permission, err error) {
	permission = &models.Permission{}

	query := repo.DB.WithContext(ctx).
		Select("id", "name", "created_at", "updated_at").
		Where("name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return permission, nil
}
