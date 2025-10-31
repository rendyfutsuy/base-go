package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

// CreateRole creates a new role information entry in the database.
//
// It takes a ToDBCreateRole parameter and returns an Role pointer and an error.
func (repo *roleRepository) CreateRole(ctx context.Context, roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error) {
	now := time.Now().UTC()

	roleRes = &models.Role{
		Name: roleReq.Name,
		Description: utils.NullString{
			String: roleReq.Description,
			Valid:  true,
		},
		CreatedAt: now,
		UpdatedAt: utils.NullTime{
			Time:  now,
			Valid: true,
		},
	}

	// Create role - GORM will insert all fields from struct
	err = repo.DB.WithContext(ctx).Create(roleRes).Error
	if err != nil {
		return nil, err
	}

	// Reload only the fields we need
	err = repo.DB.WithContext(ctx).
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("id = ?", roleRes.ID).
		First(roleRes).Error
	if err != nil {
		return nil, err
	}

	// Sync Permission Group to Role
	permissionGroupIds := dto.ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: roleReq.PermissionGroups,
	}

	// assign permission group
	err = repo.ReAssignPermissionGroup(ctx, roleRes.ID, permissionGroupIds)
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when assigning Permission Group to Role")
	}

	return roleRes, nil
}

// GetRoleByID retrieves an role information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an Role pointer and an error.
func (repo *roleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (role *models.Role, err error) {
	role = &models.Role{}

	// Get underlying SQL DB for raw query execution
	sqlDB, err := repo.DB.DB()
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Not Such Role Exist")
	}

	// Use Raw query with manual scanning for ARRAY_AGG using pq.Array
	query := `
		SELECT 
			role.id,
			role.name,
			role.created_at,
			role.updated_at,
			role.deleted_at,
			role.description,
			ARRAY_AGG(pg.name) AS permission_group_names,
			ARRAY_AGG(pg.id) AS permission_group_ids,
			ARRAY_AGG(DISTINCT pg.module) AS modules
		FROM 
			roles role
		LEFT JOIN
			modules_roles pgr
		ON
			role.id = pgr.role_id
		LEFT JOIN
			permission_groups pg
		ON
			pgr.permission_group_id = pg.id
		WHERE 
			role.id = $1 AND role.deleted_at IS NULL
		GROUP BY
			role.id, role.name
	`

	var permissionGroupNames []utils.NullString
	var permissionGroupIds []uuid.UUID
	var modules []utils.NullString

	err = sqlDB.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
		&role.Description,
		pq.Array(&permissionGroupNames),
		pq.Array(&permissionGroupIds),
		pq.Array(&modules),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.Logger.Error(err.Error())
			return nil, fmt.Errorf("Not Such Role Exist")
		}
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Not Such Role Exist")
	}

	// Assign scanned arrays
	role.PermissionGroupNames = permissionGroupNames
	role.PermissionGroupIds = permissionGroupIds
	role.Modules = modules

	// Handle NULL arrays from LEFT JOIN
	if role.PermissionGroupNames == nil {
		role.PermissionGroupNames = []utils.NullString{}
	}
	if role.PermissionGroupIds == nil {
		role.PermissionGroupIds = []uuid.UUID{}
	}
	if role.Modules == nil {
		role.Modules = []utils.NullString{}
	}

	// Fetch and assign permissions that role has
	permissions, err := repo.GetPermissionFromRoleId(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	// Fetch and assign permission groups that role has
	permissionGroups, err := repo.GetPermissionGroupFromRoleId(ctx, id)
	if err != nil {
		return nil, err
	}
	role.PermissionGroups = permissionGroups

	// get total user
	total, err := repo.GetTotalUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when fetching total user")
	}
	role.TotalUser = total

	return role, nil
}

// GetIndexRole retrieves a paginated list of role information from the database.
func (repo *roleRepository) GetIndexRole(ctx context.Context, req request.PageRequest) (roles []models.Role, total int, err error) {
	// Validate and sanitize pagination parameters
	validatedPage, validatedPerPage := request.ValidatePaginationParams(req.Page, req.PerPage, 100)
	offSet := (validatedPage - 1) * validatedPerPage
	searchQuery := req.Search

	// Get underlying SQL DB for raw query execution with ARRAY_AGG
	sqlDB, err := repo.DB.DB()
	if err != nil {
		return nil, 0, err
	}

	// Define allowed sort columns (whitelist to prevent SQL injection)
	allowedSortColumns := []string{"id", "name", "created_at", "updated_at", "deleted_at", "total_user"}

	// Validate and sanitize sort column and order
	sortBy := request.ValidateAndSanitizeSortColumn(req.SortBy, allowedSortColumns, "role.")
	if sortBy == "" {
		sortBy = "role.created_at" // Default if invalid
	}
	sortOrder := request.ValidateAndSanitizeSortOrder(req.SortOrder)

	// Build count query using GORM to leverage ApplySearchCondition
	countQuery := repo.DB.WithContext(ctx).
		Table("roles role").
		Select("DISTINCT role.id").
		Joins("LEFT JOIN modules_roles pgr ON role.id = pgr.role_id").
		Joins("LEFT JOIN permission_groups pg ON pgr.permission_group_id = pg.id").
		Where("role.deleted_at IS NULL")

	// Apply search using ApplySearchCondition helper
	countQuery = request.ApplySearchCondition(countQuery, searchQuery, []string{
		"role.name",
		"pg.module",
	})

	// Get total count
	var totalCount int64
	err = countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	total = int(totalCount)

	// Build base query for data retrieval - use raw SQL for ARRAY_AGG as GORM doesn't handle it well
	baseQuery := `
		SELECT 
			role.id,
			role.name,
			role.created_at,
			role.updated_at,
			role.deleted_at,
			(SELECT COUNT(*) FROM users WHERE role_id = role.id AND deleted_at IS NULL) AS total_user,
			ARRAY_AGG(DISTINCT pg.module) AS modules
		FROM 
			roles role
		LEFT JOIN
			modules_roles pgr
		ON
			role.id = pgr.role_id
		LEFT JOIN
			permission_groups pg
		ON
			pgr.permission_group_id = pg.id
		WHERE
			role.deleted_at IS NULL
	`

	// Build WHERE clause for GROUP BY and HAVING
	whereClause := " GROUP BY role.id, role.name, pg.module"
	args := []interface{}{}
	argIdx := 1

	// Apply search using BuildSearchConditionForRawSQL helper
	// Note: We use BuildSearchConditionForRawSQL instead of ApplySearchCondition because:
	// - This query uses raw SQL with sqlDB.QueryContext() to handle ARRAY_AGG with manual scanning via pq.Array()
	// - GORM cannot properly scan ARRAY_AGG results into []utils.NullString, so we need raw SQL
	// - ApplySearchCondition() requires *gorm.DB and returns *gorm.DB, which doesn't work with raw SQL strings
	// - BuildSearchConditionForRawSQL returns SQL clause string and args that can be used directly in raw SQL queries
	// - This maintains the same search logic as ApplySearchCondition() but for raw SQL context
	searchClause, searchArgs := request.BuildSearchConditionForRawSQL(searchQuery, []string{"role.name", "pg.module"}, argIdx, "HAVING")
	if searchClause != "" {
		whereClause += searchClause
		args = append(args, searchArgs...)
		argIdx += len(searchArgs)
	}

	// Build final query with pagination using parameter binding
	limitArgIdx := argIdx
	offsetArgIdx := argIdx + 1
	args = append(args, validatedPerPage, offSet)
	finalQuery := baseQuery + whereClause + " ORDER BY " + sortBy + " " + sortOrder + fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitArgIdx, offsetArgIdx)

	// Initialize roles slice
	roles = []models.Role{}

	// Execute query with manual scanning for arrays
	rows, err := sqlDB.QueryContext(ctx, finalQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		var modules []utils.NullString

		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
			&role.TotalUser,
			pq.Array(&modules),
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle NULL arrays
		if modules == nil {
			role.Modules = []utils.NullString{}
		} else {
			role.Modules = modules
		}

		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetAllRole retrieves all role information entries from the database.
//
// Returns a slice of models.Role and an error.
func (repo *roleRepository) GetAllRole(ctx context.Context) (roles []models.Role, err error) {
	err = repo.DB.WithContext(ctx).
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("deleted_at IS NULL").
		Find(&roles).Error

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// UpdateRole updates an existing role information entry in the database.
func (repo *roleRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error) {
	updates := map[string]interface{}{
		"name":        roleReq.Name,
		"description": roleReq.Description,
		"updated_at":  time.Now().UTC(),
	}

	roleRes = &models.Role{}
	err = repo.DB.WithContext(ctx).
		Model(&models.Role{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		First(roleRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role role with id %s not found", id)
		}
		return nil, err
	}

	// Sync Permission Group to Role
	permissionGroupIds := dto.ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: roleReq.PermissionGroups,
	}

	// assign permission group
	err = repo.ReAssignPermissionGroup(ctx, roleRes.ID, permissionGroupIds)
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when assigning Permission Group to Role")
	}

	return roleRes, nil
}

// SoftDeleteRole soft deletes an role role entry in the database.
func (repo *roleRepository) SoftDeleteRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error) {
	roleRes = &models.Role{}

	// GORM soft delete automatically sets deleted_at
	err = repo.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&models.Role{}).Error

	if err != nil {
		return nil, err
	}

	// Get the deleted role (with Unscoped to include soft deleted)
	err = repo.DB.WithContext(ctx).
		Unscoped().
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("id = ?", id).
		First(roleRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role role with id %s not found", id)
		}
		return nil, err
	}

	return roleRes, nil
}

// CountRole retrieves the count of role information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountRole(ctx context.Context) (count *int, err error) {
	var result int64
	err = repo.DB.WithContext(ctx).
		Model(&models.Role{}).
		Count(&result).Error

	if err != nil {
		return nil, err
	}

	resultInt := int(result)
	count = &resultInt
	return count, nil
}

// RoleNameIsNotDuplicated checks if the provided role name is not duplicated in the database.
func (repo *roleRepository) RoleNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.Role{}).
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

// GetDuplicatedRole retrieves the role information with the given name and excluded ID from the database.
func (repo *roleRepository) GetDuplicatedRole(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	role = &models.Role{}

	query := repo.DB.WithContext(ctx).
		Select("id", "name", "created_at", "updated_at").
		Where("name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// RoleNameIsNotDuplicatedOnSoftDeleted checks if the provided role name is not duplicated in the database.
func (repo *roleRepository) RoleNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.Role{}).
		Unscoped().
		Where("name = ?", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetDuplicatedRoleOnSoftDeleted retrieves the role information with the given name and excluded ID from the database.
func (repo *roleRepository) GetDuplicatedRoleOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	role = &models.Role{}

	query := repo.DB.WithContext(ctx).
		Unscoped().
		Select("id", "name", "created_at", "updated_at").
		Where("name = ?", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}
