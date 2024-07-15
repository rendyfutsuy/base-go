package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	class "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type classRepository struct {
	Conn *sql.DB
}

func NewClassRepository(Conn *sql.DB) class.Repository {
	return &classRepository{Conn}
}

func (repo *classRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *classRepository) CreateClass(classReq dto.ToDBCreateClass) (classRes *models.Class, err error) {

	classRes = new(models.Class)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO classes
			(name, code, created_at, created_by)
		VALUES
			($1, $2, $3, $4)
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		classReq.Name,
		classReq.Code,
		createdAtString,
		classReq.CreatedByID,
	).Scan(
		&classRes.ID,
		&classRes.Name,
		&classRes.Code,
		&classRes.CreatedAt,
		&classRes.CreatedByID,
		&classRes.UpdatedAt,
		&classRes.UpdatedByID,
		&classRes.DeletedAt,
		&classRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return classRes, err
}

func (repo *classRepository) GetClassByID(id uuid.UUID) (class *models.Class, err error) {
	class = new(models.Class)
	err = repo.Conn.QueryRow(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			classes 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&class.ID,
		&class.Name,
		&class.Code,
		&class.CreatedAt,
		&class.CreatedByID,
		&class.UpdatedAt,
		&class.UpdatedByID,
		&class.DeletedAt,
		&class.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("class with id %s not found", id)
		}

		return nil, err
	}

	return class, err
}

func (repo *classRepository) GetIndexClass(req request.PageRequest) (classes []models.Class, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM classes"
	countQuery := "SELECT COUNT(*) FROM classes"
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
		var class models.Class
		err = rows.Scan(
			&class.ID,
			&class.Name,
			&class.Code,
			&class.CreatedAt,
			&class.CreatedByID,
			&class.UpdatedAt,
			&class.UpdatedByID,
			&class.DeletedAt,
			&class.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		classes = append(classes, class)
	}

	return classes, total, err
}

func (repo *classRepository) GetAllClass() (classes []models.Class, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			classes
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var class models.Class
		err = rows.Scan(
			&class.ID,
			&class.Name,
			&class.Code,
			&class.CreatedAt,
			&class.CreatedByID,
			&class.UpdatedAt,
			&class.UpdatedByID,
			&class.DeletedAt,
			&class.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		classes = append(classes, class)
	}

	return classes, err
}

func (repo *classRepository) UpdateClass(id uuid.UUID, classReq dto.ToDBUpdateClass) (classRes *models.Class, err error) {

	classRes = new(models.Class)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE classes SET 
			name = $1, updated_at = $2, updated_by = $3
		WHERE 
			id = $4 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		classReq.Name,
		updatedAtString,
		classReq.UpdatedByID,
		id,
	).Scan(
		&classRes.ID,
		&classRes.Name,
		&classRes.Code,
		&classRes.CreatedAt,
		&classRes.CreatedByID,
		&classRes.UpdatedAt,
		&classRes.UpdatedByID,
		&classRes.DeletedAt,
		&classRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("class with id %s not found", id)
		}

		return nil, err
	}

	return classRes, err
}

func (repo *classRepository) SoftDeleteClass(id uuid.UUID, classReq dto.ToDBDeleteClass) (classRes *models.Class, err error) {

	classRes = new(models.Class)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE classes SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		deletedAtString,
		classReq.DeletedByID,
		id,
	).Scan(
		&classRes.ID,
		&classRes.Name,
		&classRes.Code,
		&classRes.CreatedAt,
		&classRes.CreatedByID,
		&classRes.UpdatedAt,
		&classRes.UpdatedByID,
		&classRes.DeletedAt,
		&classRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("class with id %s not found", id)
		}

		return nil, err
	}

	return classRes, err
}

func (repo *classRepository) CountClass() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			classes`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
