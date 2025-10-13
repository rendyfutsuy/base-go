package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
)

// GetPermissionByID retrieves an permission information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an Permission pointer and an error.
func (repo *roleRepository) GetPermissionByID(id uuid.UUID) (permission *models.Permission, err error) {
	// initialize permission variable
	permission = new(models.Permission)

	// fetch data from database by id that passed
	// assign return value to permission variable
	err = repo.Conn.QueryRow(
		`SELECT 
			permission.id,
			permission.name,
			permission.created_at,
			permission.updated_at,
			permission.deleted_at
		FROM 
			permissions permission
		WHERE 
			permission.id = $1 AND permission.deleted_at IS NULL`,
		id,
	).Scan(
		&permission.ID,
		&permission.Name,
		&permission.CreatedAt,
		&permission.UpdatedAt,
		&permission.DeletedAt,
	)

	// if error occurs, return error
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission permission with id %s not found", id)
		}

		return nil, err
	}

	return permission, err
}

// GetIndexPermission retrieves a paginated list of permission information from the database.
//
// It takes a PageRequest parameter and returns a slice of Permission, the total number of
// permission information entries, and an error.
// its can search by permission name, permission code, permission alias_1, permission alias_2, permission alias_3, permission alias_4, permission address, permission email, permission phone_number, type name
func (repo *roleRepository) GetIndexPermission(req request.PageRequest) (permissions []models.Permission, total int, err error) {

	// initialize: pagination page, search query to local variable
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query

	// construct select fillable
	selectFill := `
		permission.id,
		permission.name,
		permission.created_at,
		permission.updated_at,
		permission.deleted_at
	`

	// append select and join query to base query
	baseQuery := "SELECT " + selectFill + " FROM permissions permission "

	// append join query to count query
	countQuery := "SELECT COUNT(*) FROM permissions permission "

	// initialize common query for deleted condition
	whereClause := " WHERE permission.deleted_at IS NULL"

	// assign search query, based on search able field.
	if searchQuery != "" {
		searchPermission := "permission.name ILIKE '%' || $1 || '%'"
		whereClause += " AND (" + searchPermission + ")"
	}

	// Initialize Default sorting
	sortBy := "permission.created_at"
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
		err = repo.Conn.QueryRow(countQuery+whereClause, searchQuery).Scan(&total)
	} else {
		err = repo.Conn.QueryRow(countQuery + whereClause).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	// retrieve paginated
	rows := new(sql.Rows)
	if searchQuery != "" {
		rows, err = repo.Conn.Query(baseQuery+whereClause+orderClause+limitClause, searchQuery)
	} else {
		rows, err = repo.Conn.Query(baseQuery + whereClause + orderClause + limitClause)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// assign pagination as models.Permission
	for rows.Next() {
		var permission models.Permission
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, total, err
}

// GetAllPermission retrieves all permission information entries from the database.
//
// Returns a slice of models.Permission and an error.
func (repo *roleRepository) GetAllPermission() (permissions []models.Permission, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			permission.id,
			permission.name,
			permission.created_at,
			permission.updated_at,
			permission.deleted_at
		FROM 
			permissions permission
		WHERE
			permission.deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var permission models.Permission
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&permission.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, err
}

// CountPermission retrieves the count of permission information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *roleRepository) CountPermission() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			permissions`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}

// PermissionNameIsNotDuplicated checks if the provided permission name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *roleRepository) PermissionNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			permissions
		WHERE
			name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	result := 0

	// assert name is nt duplicated
	err := repo.Conn.QueryRow(baseQuery, params...).Scan(&result)

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

// GetDuplicatedPermission retrieves the permission information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the permission information to retrieve.
// - excludedId: the ID of the permission information to exclude from the result.
//
// Returns:
// - permission: a pointer to the retrieved permission information.
// - err: an error if there was a problem retrieving the permission information.
func (repo *roleRepository) GetDuplicatedPermission(name string, excludedId uuid.UUID) (permission *models.Permission, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			permissions
		WHERE
			name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize permission
	permission = &models.Permission{}

	// assert name is not duplicated
	err = repo.Conn.QueryRow(baseQuery, params...).Scan(
		&permission.ID,
		&permission.Name,
		&permission.CreatedAt,
		&permission.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated Permission
	return permission, nil
}
