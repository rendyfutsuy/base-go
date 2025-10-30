package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"context"
)

// GetPermissionGroupByID retrieves an permission_group information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an PermissionGroup pointer and an error.
func (repo *roleRepository) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	// initialize permission_group variable
	permissionGroup = new(models.PermissionGroup)

	query := `
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
            pg.id = $1 AND pg.deleted_at IS NULL
        GROUP BY
            pg.id, pg.name;
    `

	// Fetch data from the database by id that was passed
	// Assign return value to permission_group variable
	err = repo.Conn.QueryRowContext(ctx, query, id).Scan(
		&permissionGroup.ID,
		&permissionGroup.Name,
		pq.Array(&permissionGroup.PermissionNames),
		&permissionGroup.CreatedAt,
		&permissionGroup.UpdatedAt,
		&permissionGroup.DeletedAt,
	)

	// if error occurs, return error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission_group permission_group with id %s not found", id)
		}

		return nil, err
	}

	return permissionGroup, err
}

// GetIndexPermissionGroup retrieves a paginated list of permission_group information from the database.
//
// It takes a PageRequest parameter and returns a slice of PermissionGroup, the total number of
// permission_group information entries, and an error.
// its can search by permission_group name, permission_group code, permission_group alias_1, permission_group alias_2, permission_group alias_3, permission_group alias_4, permission_group address, permission_group email, permission_group phone_number, type name
func (repo *roleRepository) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (permissionGroups []models.PermissionGroup, total int, err error) {

	// initialize: pagination page, search query to local variable
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query

	// construct select fillable
	selectFill := `
		permission_group.id,
		permission_group.name,
		permission_group.created_at,
		permission_group.updated_at,
		permission_group.deleted_at
	`

	// append select and join query to base query
	baseQuery := "SELECT " + selectFill + " FROM permission_groups permission_group "

	// append join query to count query
	countQuery := "SELECT COUNT(*) FROM permission_groups permission_group "

	// initialize common query for deleted condition
	whereClause := " WHERE permission_group.deleted_at IS NULL"

	// assign search query, based on search able field.
	if searchQuery != "" {
		searchPermissionGroup := "permission_group.name ILIKE '%' || $1 || '%'"
		whereClause += " AND (" + searchPermissionGroup + ")"
	}

	// Initialize Default sorting
	sortBy := "permission_group.created_at"
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
		err = repo.Conn.QueryRowContext(ctx, countQuery + whereClause).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	// retrieve paginated
	rows := new(sql.Rows)
	if searchQuery != "" {
		rows, err = repo.Conn.QueryContext(ctx, baseQuery+whereClause+orderClause+limitClause, searchQuery)
	} else {
		rows, err = repo.Conn.QueryContext(ctx, baseQuery + whereClause + orderClause + limitClause)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// assign pagination as models.PermissionGroup
	for rows.Next() {
		var permissionGroup models.PermissionGroup
		err = rows.Scan(
			&permissionGroup.ID,
			&permissionGroup.Name,
			&permissionGroup.CreatedAt,
			&permissionGroup.UpdatedAt,
			&permissionGroup.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		permissionGroups = append(permissionGroups, permissionGroup)
	}

	return permissionGroups, total, err
}

// GetAllPermissionGroup retrieves all permission_group information entries from the database.
//
// Returns a slice of models.PermissionGroup and an error.
func (repo *roleRepository) GetAllPermissionGroup(ctx context.Context) (permissionGroups []models.PermissionGroup, err error) {
	rows, err := repo.Conn.QueryContext(ctx, 
		`SELECT 
			permission_group.id,
			permission_group.name,
			permission_group.module,
			permission_group.created_at,
			permission_group.updated_at,
			permission_group.deleted_at
		FROM 
			permission_groups permission_group
		WHERE
			permission_group.deleted_at IS NULL
		ORDER BY permission_group.module ASC`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permissionGroup models.PermissionGroup
		err = rows.Scan(
			&permissionGroup.ID,
			&permissionGroup.Name,
			&permissionGroup.Module,
			&permissionGroup.CreatedAt,
			&permissionGroup.UpdatedAt,
			&permissionGroup.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		permissionGroups = append(permissionGroups, permissionGroup)
	}

	return permissionGroups, err
}

// CountPermissionGroup retrieves the count of permission_group information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountPermissionGroup(ctx context.Context) (count *int, err error) {
	err = repo.Conn.QueryRowContext(ctx, 
		`SELECT 
			COUNT(*)
		FROM 
			permission_groups`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}

// PermissionGroupNameIsNotDuplicated checks if the provided permission_group name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *roleRepository) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			permission_groups
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

// GetDuplicatedPermissionGroup retrieves the permission_group information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the permission_group information to retrieve.
// - excludedId: the ID of the permission_group information to exclude from the result.
//
// Returns:
// - permission_group: a pointer to the retrieved permission_group information.
// - err: an error if there was a problem retrieving the permission_group information.
func (repo *roleRepository) GetDuplicatedPermissionGroup(ctx context.Context, name string, excludedId uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			permission_groups
		WHERE
			name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize permission_group
	permissionGroup = &models.PermissionGroup{}

	// assert name is not duplicated
	err = repo.Conn.QueryRowContext(ctx, baseQuery, params...).Scan(
		&permissionGroup.ID,
		&permissionGroup.Name,
		&permissionGroup.CreatedAt,
		&permissionGroup.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated PermissionGroup
	return permissionGroup, nil
}
