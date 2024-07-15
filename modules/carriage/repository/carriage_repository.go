package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	carriage "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type carriageRepository struct {
	Conn *sql.DB
}

func NewCarriageRepository(Conn *sql.DB) carriage.Repository {
	return &carriageRepository{Conn}
}

func (repo *carriageRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *carriageRepository) CreateCarriage(carriageReq dto.ToDBCreateCarriage) (carriageRes *models.Carriage, err error) {

	carriageRes = new(models.Carriage)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO carriages
			(name, code, created_at, created_by)
		VALUES
			($1, $2, $3, $4)
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		carriageReq.Name,
		carriageReq.Code,
		createdAtString,
		carriageReq.CreatedByID,
	).Scan(
		&carriageRes.ID,
		&carriageRes.Name,
		&carriageRes.Code,
		&carriageRes.CreatedAt,
		&carriageRes.CreatedByID,
		&carriageRes.UpdatedAt,
		&carriageRes.UpdatedByID,
		&carriageRes.DeletedAt,
		&carriageRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return carriageRes, err
}

func (repo *carriageRepository) GetCarriageByID(id uuid.UUID) (carriage *models.Carriage, err error) {
	carriage = new(models.Carriage)
	err = repo.Conn.QueryRow(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			carriages 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&carriage.ID,
		&carriage.Name,
		&carriage.Code,
		&carriage.CreatedAt,
		&carriage.CreatedByID,
		&carriage.UpdatedAt,
		&carriage.UpdatedByID,
		&carriage.DeletedAt,
		&carriage.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("carriage with id %s not found", id)
		}

		return nil, err
	}

	return carriage, err
}

func (repo *carriageRepository) GetIndexCarriage(req request.PageRequest) (carriages []models.Carriage, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM carriages"
	countQuery := "SELECT COUNT(*) FROM carriages"
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
		var carriage models.Carriage
		err = rows.Scan(
			&carriage.ID,
			&carriage.Name,
			&carriage.Code,
			&carriage.CreatedAt,
			&carriage.CreatedByID,
			&carriage.UpdatedAt,
			&carriage.UpdatedByID,
			&carriage.DeletedAt,
			&carriage.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		carriages = append(carriages, carriage)
	}

	return carriages, total, err
}

func (repo *carriageRepository) GetAllCarriage() (carriages []models.Carriage, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			carriages
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var carriage models.Carriage
		err = rows.Scan(
			&carriage.ID,
			&carriage.Name,
			&carriage.Code,
			&carriage.CreatedAt,
			&carriage.CreatedByID,
			&carriage.UpdatedAt,
			&carriage.UpdatedByID,
			&carriage.DeletedAt,
			&carriage.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		carriages = append(carriages, carriage)
	}

	return carriages, err
}

func (repo *carriageRepository) UpdateCarriage(id uuid.UUID, carriageReq dto.ToDBUpdateCarriage) (carriageRes *models.Carriage, err error) {

	carriageRes = new(models.Carriage)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE carriages SET 
			name = $1, updated_at = $2, updated_by = $3
		WHERE 
			id = $4 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		carriageReq.Name,
		updatedAtString,
		carriageReq.UpdatedByID,
		id,
	).Scan(
		&carriageRes.ID,
		&carriageRes.Name,
		&carriageRes.Code,
		&carriageRes.CreatedAt,
		&carriageRes.CreatedByID,
		&carriageRes.UpdatedAt,
		&carriageRes.UpdatedByID,
		&carriageRes.DeletedAt,
		&carriageRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("carriage with id %s not found", id)
		}

		return nil, err
	}

	return carriageRes, err
}

func (repo *carriageRepository) SoftDeleteCarriage(id uuid.UUID, carriageReq dto.ToDBDeleteCarriage) (carriageRes *models.Carriage, err error) {

	carriageRes = new(models.Carriage)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE carriages SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		deletedAtString,
		carriageReq.DeletedByID,
		id,
	).Scan(
		&carriageRes.ID,
		&carriageRes.Name,
		&carriageRes.Code,
		&carriageRes.CreatedAt,
		&carriageRes.CreatedByID,
		&carriageRes.UpdatedAt,
		&carriageRes.UpdatedByID,
		&carriageRes.DeletedAt,
		&carriageRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("carriage with id %s not found", id)
		}

		return nil, err
	}

	return carriageRes, err
}

func (repo *carriageRepository) CountCarriage() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			carriages`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
