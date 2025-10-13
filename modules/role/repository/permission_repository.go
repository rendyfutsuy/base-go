package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	"github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/modules/role/dto"
	"github.com/rendyfutsuy/base-go.git/utils"
)

func (repo *roleRepository) CreatePermission(permissionReq dto.ToDBCreatePermission) (permissionRes *models.Permission, err error) {

	permissionRes = new(models.Permission)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO permissions
			(name, created_at)
		VALUES
			($1, $2, $3, $4)
		RETURNING 
			id, name, created_at, updated_at`,
		permissionReq.Name,
		createdAtString,
		permissionReq.CreatedByID,
	).Scan(
		&permissionRes.ID,
		&permissionRes.Name,
		&permissionRes.CreatedAt,
		&permissionRes.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return permissionRes, err
}

func (repo *roleRepository) GetPermissionByID(id uuid.UUID) (permission *models.Permission, err error) {
	permission = new(models.Permission)
	err = repo.Conn.QueryRow(
		`SELECT 
			id, name, created_at, updated_at
		FROM 
			permissions 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&permission.ID,
		&permission.Name,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission with id %s not found", id)
		}

		return nil, err
	}

	return permission, err
}

func (repo *roleRepository) GetIndexPermission(req request.PageRequest) (permissions []models.Permission, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM permissions"
	countQuery := "SELECT COUNT(*) FROM permissions"
	whereClause := " WHERE deleted_at IS NULL"
	if searchQuery != "" {
		whereClause += " AND (name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%')"
	}

	// Default sorting
	sortBy := "created_at"
	sortOrder := "DESC" // Sort from newest to oldest
	if req.SortBy != "" {
		sortBy = req.SortBy
		if req.SortOrder != "" {
			sortOrder = req.SortOrder
		}
	}

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

	for rows.Next() {
		var permission models.Permission
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, total, err
}

func (repo *roleRepository) GetAllPermission() (permissions []models.Permission, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id, name, code, created_at, updated_at
		FROM 
			permissions
		WHERE
			deleted_at IS NULL`,
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
		)

		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	return permissions, err
}

func (repo *roleRepository) UpdatePermission(id uuid.UUID, permissionReq dto.ToDBUpdatePermission) (permissionRes *models.Permission, err error) {

	permissionRes = new(models.Permission)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE permissions SET 
			name = $1, updated_at = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at`,
		permissionReq.Name,
		updatedAtString,
		id,
	).Scan(
		&permissionRes.ID,
		&permissionRes.Name,
		&permissionRes.CreatedAt,
		&permissionRes.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission with id %s not found", id)
		}

		return nil, err
	}

	return permissionRes, err
}

func (repo *roleRepository) SoftDeletePermission(id uuid.UUID, permissionReq dto.ToDBDeletePermission) (permissionRes *models.Permission, err error) {

	permissionRes = new(models.Permission)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE permissions SET 
			deleted_at = $1
		WHERE 
			id = $2 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at`,
		deletedAtString,
		id,
	).Scan(
		&permissionRes.ID,
		&permissionRes.Name,
		&permissionRes.CreatedAt,
		&permissionRes.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission with id %s not found", id)
		}

		return nil, err
	}

	return permissionRes, err
}

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
