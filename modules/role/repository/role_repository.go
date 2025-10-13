package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	"github.com/rendyfutsuy/base-go.git/models"
	role "github.com/rendyfutsuy/base-go.git/modules/role"
	"github.com/rendyfutsuy/base-go.git/modules/role/dto"
	"github.com/rendyfutsuy/base-go.git/utils"
)

type roleRepository struct {
	Conn *sql.DB
}

func NewRoleRepository(Conn *sql.DB) role.Repository {
	return &roleRepository{Conn}
}

func (repo *roleRepository) CreateTable(sqlFilePath string) (err error) {

	sqlCommands, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return err
	}

	_, err = repo.Conn.Exec(string(sqlCommands))
	if err != nil {
		return err
	}

	return err
}

func (repo *roleRepository) CreateRole(roleReq dto.ToDBCreateRole) (roleRes *models.Role, err error) {

	roleRes = new(models.Role)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO roles
			(name, created_at)
		VALUES
			($1, $2, $3, $4)
		RETURNING 
			id, name, created_at, updated_at`,
		roleReq.Name,
		createdAtString,
		roleReq.CreatedByID,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return roleRes, err
}

func (repo *roleRepository) GetRoleByID(id uuid.UUID) (role *models.Role, err error) {
	role = new(models.Role)
	err = repo.Conn.QueryRow(
		`SELECT 
			id, name, created_at, updated_at
		FROM 
			roles 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role with id %s not found", id)
		}

		return nil, err
	}

	return role, err
}

func (repo *roleRepository) GetIndexRole(req request.PageRequest) (roles []models.Role, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM roles"
	countQuery := "SELECT COUNT(*) FROM roles"
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
		var role models.Role
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		roles = append(roles, role)
	}

	return roles, total, err
}

func (repo *roleRepository) GetAllRole() (roles []models.Role, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id, name, code, created_at, updated_at
		FROM 
			roles
		WHERE
			deleted_at IS NULL`,
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
		)

		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, err
}

func (repo *roleRepository) UpdateRole(id uuid.UUID, roleReq dto.ToDBUpdateRole) (roleRes *models.Role, err error) {

	roleRes = new(models.Role)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE roles SET 
			name = $1, updated_at = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at`,
		roleReq.Name,
		updatedAtString,
		id,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role with id %s not found", id)
		}

		return nil, err
	}

	return roleRes, err
}

func (repo *roleRepository) SoftDeleteRole(id uuid.UUID, roleReq dto.ToDBDeleteRole) (roleRes *models.Role, err error) {

	roleRes = new(models.Role)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE roles SET 
			deleted_at = $1
		WHERE 
			id = $2 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at`,
		deletedAtString,
		id,
	).Scan(
		&roleRes.ID,
		&roleRes.Name,
		&roleRes.CreatedAt,
		&roleRes.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role with id %s not found", id)
		}

		return nil, err
	}

	return roleRes, err
}

func (repo *roleRepository) CountRole() (count *int, err error) {
	err = repo.Conn.QueryRow(
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
