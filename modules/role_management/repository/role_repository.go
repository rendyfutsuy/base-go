package repository

import (
	"database/sql"
	"fmt"
	"time"

	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

// CreateRole creates a new role information entry in the database.
//
// It takes a ToDBCreateRole parameter and returns an Role pointer and an error.
func (repo *roleRepository) CreateRole(ctx context.Context, roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error) {

	// initialize: role role model, time format to created at string,
	roleRes = new(models.Role)
	timeFormat := constants.FormatTimezone
	createdAtString := time.Now().UTC().Format(timeFormat)

	// execute query to insert role role
	// assign return value to roleRes variable
	err = repo.Conn.QueryRowContext(ctx,
		`INSERT INTO roles
			(name, created_at, updated_at, description, deletable)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING 
			id, name, created_at, updated_at, deleted_at`,
		roleReq.Name,
		createdAtString,
		createdAtString,
		roleReq.Description,
		true,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
		&roleRes.DeletedAt,
	)
	// if error occurs, return error
	if err != nil {
		return nil, err
	}

	// Sync Permission Group to Role
	//mapping permission group
	permissionGroupIds := dto.ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: roleReq.PermissionGroups,
	}

	// assign permission group
	err = repo.ReAssignPermissionGroup(ctx, roleRes.ID, permissionGroupIds)

	// if error occurs, return error
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when assigning Permission Group to Role")
	}

	return roleRes, err
}

// GetRoleByID retrieves an role information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an Role pointer and an error.
func (repo *roleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (role *models.Role, err error) {
	// initialize role variable
	role = new(models.Role)

	// fetch data from database by id that passed
	// assign return value to role variable
	err = repo.Conn.QueryRowContext(ctx,
		`SELECT 
			role.id,
			role.name,
			role.created_at,
			role.updated_at,
			role.deleted_at,
			role.description,
			ARRAY_AGG(pg.name) AS permission_groups,
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
			role.id, role.name`,
		id,
	).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.DeletedAt,
		&role.Description,
		pq.Array(&role.PermissionGroupNames),
		pq.Array(&role.PermissionGroupIds),
		pq.Array(&role.Modules),
	)

	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, fmt.Errorf("Not Such Role Exist")
	}

	// Fetch and assign permissions that role has
	permissions, err := repo.GetPermissionFromRoleId(ctx, id)

	if err != nil {
		return nil, err
	}

	// assign permissions to Role model
	role.Permissions = permissions

	// Fetch and assign permission groups that role has
	permissionGroups, err := repo.GetPermissionGroupFromRoleId(ctx, id)

	if err != nil {
		return nil, err
	}

	// assign permission groups to Role model
	role.PermissionGroups = permissionGroups

	// if error occurs, return error
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when fetching permission group..")
	}
	// get total user
	total, err := repo.GetTotalUser(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("Something Wrong when fetching total user")
	}

	// assign total user to Role model
	role.TotalUser = total

	return role, err
}

// GetIndexRole retrieves a paginated list of role information from the database.
//
// It takes a PageRequest parameter and returns a slice of Role, the total number of
// role information entries, and an error.
// its can search by role name, role code, role alias_1, role alias_2, role alias_3, role alias_4, role address, role email, role phone_number, type name
func (repo *roleRepository) GetIndexRole(ctx context.Context, req request.PageRequest) (roles []models.Role, total int, err error) {
	// initialize: pagination page, search query to local variable
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query

	// Join Query
	JoinFill := ` LEFT JOIN
			modules_roles pgr
		ON
			role.id = pgr.role_id
		LEFT JOIN
			permission_groups pg
		ON
			pgr.permission_group_id = pg.id`

	GroupBy := ` GROUP BY
			role.id, role.name`

	// construct select fillable
	selectFill := `
		role.id,
		role.name,
		role.created_at,
		role.updated_at,
		role.deleted_at,
		(SELECT COUNT(*) FROM users WHERE role_id = role.id AND deleted_at IS NULL) AS total_user,
		ARRAY_AGG(DISTINCT pg.module) AS modules
	`

	// append select and join query to base query
	baseQuery := "SELECT " + selectFill + " FROM roles role" + JoinFill

	// append join query to count query
	countQuery := "SELECT COUNT(DISTINCT role.id) FROM roles role" + JoinFill

	// initialize common query for deleted condition
	whereClause := " WHERE role.deleted_at IS NULL"

	// assign search query, based on searchable field.
	if searchQuery != "" {
		searchModule := "pg.module ILIKE '%' || $1 || '%'"
		searchRole := "role.name ILIKE '%' || $1 || '%'"
		whereClause += " AND (" + searchRole + " OR " + searchModule + ")"
	}

	// Initialize Default sorting
	sortBy := "role.created_at"
	sortOrder := "DESC" // Sort from newest to oldest
	if req.SortBy != "" {
		sortBy = req.SortBy
		if req.SortOrder != "" {
			sortOrder = req.SortOrder
		}
	}

	// initialize Default Sorting Query
	orderClause := " ORDER BY " + sortBy + " " + sortOrder
	limitClause := fmt.Sprintf(" LIMIT %d OFFSET %d", req.PerPage, offSet)

	// count total
	if searchQuery != "" {
		err = repo.Conn.QueryRowContext(ctx, countQuery+whereClause, searchQuery).Scan(&total)
	} else {
		err = repo.Conn.QueryRowContext(ctx, countQuery+whereClause).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	// retrieve paginated
	rows := new(sql.Rows)
	if searchQuery != "" {
		rows, err = repo.Conn.QueryContext(ctx, baseQuery+whereClause+GroupBy+orderClause+limitClause, searchQuery)
	} else {
		rows, err = repo.Conn.QueryContext(ctx, baseQuery+whereClause+GroupBy+orderClause+limitClause)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// assign pagination as models.Role
	for rows.Next() {
		var role models.Role
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
			&role.TotalUser,
			pq.Array(&role.Modules),
		)

		if err != nil {
			return nil, 0, err
		}

		roles = append(roles, role)
	}

	return roles, total, err
}

// GetAllRole retrieves all role information entries from the database.
//
// Returns a slice of models.Role and an error.
func (repo *roleRepository) GetAllRole(ctx context.Context) (roles []models.Role, err error) {
	rows, err := repo.Conn.QueryContext(ctx,
		`SELECT 
			role.id,
			role.name,
			role.created_at,
			role.updated_at,
			role.deleted_at
		FROM 
			roles role
		WHERE
			role.deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
			&role.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, err
}

// UpdateRole updates an existing role information entry in the database.
//
// It takes an ID of the role information and a ToDBUpdateRole parameter.
// It returns an Role pointer and an error.
//
// The function updates the role information in the database with the provided ID.
// It sets the name, updated_at, updated_by, email, phone_number, role_type_id,
// alias_1, alias_2, alias_3, alias_4, and address fields of the role information.
// If the role information with the provided ID is not found, it returns an error.
// If there is an error during the update, it returns the error.
func (repo *roleRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error) {
	roleRes = new(models.Role)
	timeFormat := constants.FormatTimezone
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRowContext(ctx,
		`UPDATE roles SET 
			name = $1,
			updated_at = $2,
			description = $4
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id,
			name,
			created_at,
			updated_at,
			deleted_at`,
		roleReq.Name,
		updatedAtString,
		id,
		roleReq.Description,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
		&roleRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role role with id %s not found", id)
		}

		return nil, err
	}

	// Sync Permission Group to Role
	//mapping permission group
	permissionGroupIds := dto.ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: roleReq.PermissionGroups,
	}

	// assign permission group
	err = repo.ReAssignPermissionGroup(ctx, roleRes.ID, permissionGroupIds)

	// if error occurs, return error
	if err != nil {
		return nil, fmt.Errorf("Something Wrong when assigning Permission Group to Role")
	}

	return roleRes, err
}

// SoftDeleteRole soft deletes an role role entry in the database.
//
// It takes an id of type uuid.UUID and an roleReq of type dto.ToDBDeleteRole as parameters.
// It returns the soft deleted role role entry of type models.Role and an error.
func (repo *roleRepository) SoftDeleteRole(ctx context.Context, id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error) {

	roleRes = new(models.Role)
	timeFormat := constants.FormatTimezone
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRowContext(ctx,
		`UPDATE roles SET 
			deleted_at = $1
		WHERE 
			id = $2 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at, deleted_at`,
		deletedAtString,
		id,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
		&roleRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role role with id %s not found", id)
		}

		return nil, err
	}

	return roleRes, err
}

// CountRole retrieves the count of role information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountRole(ctx context.Context) (count *int, err error) {
	err = repo.Conn.QueryRowContext(ctx,
		`SELECT 
			COUNT(*)
		FROM 
			roles`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}

// RoleNameIsNotDuplicated checks if the provided role name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *roleRepository) RoleNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			roles
		WHERE
			name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	result := 0

	// assert name is nt duplicated
	err := repo.Conn.QueryRowContext(ctx, baseQuery, params...).Scan(&result)

	// if have error, return false and error
	if err != nil {
		return false, err
	}

	// if duplicated name, return false
	if result > 0 {
		return false, nil
	}

	// if not duplicated name, return true
	return true, err
}

// GetDuplicatedRole retrieves the role information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the role information to retrieve.
// - excludedId: the ID of the role information to exclude from the result.
//
// Returns:
// - role: a pointer to the retrieved role information.
// - err: an error if there was a problem retrieving the role information.
func (repo *roleRepository) GetDuplicatedRole(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			roles
		WHERE
			name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize role
	role = &models.Role{}

	// assert name is not duplicated
	err = repo.Conn.QueryRowContext(ctx, baseQuery, params...).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated Role
	return role, nil
}

// RoleNameIsNotDuplicatedOnSoftDeleted checks if the provided role name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *roleRepository) RoleNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			roles
		WHERE
			name = $1`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	result := 0

	// assert name is nt duplicated
	err := repo.Conn.QueryRowContext(ctx, baseQuery, params...).Scan(&result)

	// if have error, return false and error
	if err != nil {
		return false, err
	}

	// if duplicated name, return false
	if result > 0 {
		return false, nil
	}

	// if not duplicated name, return true
	return true, err
}

// GetDuplicatedRoleOnSoftDeleted retrieves the role information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the role information to retrieve.
// - excludedId: the ID of the role information to exclude from the result.
//
// Returns:
// - role: a pointer to the retrieved role information.
// - err: an error if there was a problem retrieving the role information.
func (repo *roleRepository) GetDuplicatedRoleOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			roles
		WHERE
			name = $1`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize role
	role = &models.Role{}

	// assert name is not duplicated
	err = repo.Conn.QueryRowContext(ctx, baseQuery, params...).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated Role
	return role, nil
}
