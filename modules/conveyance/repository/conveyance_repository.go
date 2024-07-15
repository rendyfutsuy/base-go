package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	conveyance "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type conveyanceRepository struct {
	Conn *sql.DB
}

func NewConveyanceRepository(Conn *sql.DB) conveyance.Repository {
	return &conveyanceRepository{Conn}
}

func (repo *conveyanceRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *conveyanceRepository) CreateConveyance(conveyanceReq dto.ToDBCreateConveyance) (conveyanceRes *models.Conveyance, err error) {

	conveyanceRes = new(models.Conveyance)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO conveyances
			(name, code, type, created_at, created_by)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING 
			* `,
		conveyanceReq.Name,
		conveyanceReq.Code,
		conveyanceReq.Type,
		createdAtString,
		conveyanceReq.CreatedByID,
	).Scan(
		&conveyanceRes.ID,
		&conveyanceRes.Name,
		&conveyanceRes.Code,
		&conveyanceRes.Type,
		&conveyanceRes.CreatedAt,
		&conveyanceRes.CreatedByID,
		&conveyanceRes.UpdatedAt,
		&conveyanceRes.UpdatedByID,
		&conveyanceRes.DeletedAt,
		&conveyanceRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return conveyanceRes, err
}

func (repo *conveyanceRepository) GetConveyanceByID(id uuid.UUID) (conveyance *models.Conveyance, err error) {
	conveyance = new(models.Conveyance)
	err = repo.Conn.QueryRow(
		`SELECT 
			*
		FROM 
			conveyances 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&conveyance.ID,
		&conveyance.Name,
		&conveyance.Code,
		&conveyance.Type,
		&conveyance.CreatedAt,
		&conveyance.CreatedByID,
		&conveyance.UpdatedAt,
		&conveyance.UpdatedByID,
		&conveyance.DeletedAt,
		&conveyance.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conveyance with id %s not found", id)
		}

		return nil, err
	}

	return conveyance, err
}

func (repo *conveyanceRepository) GetIndexConveyance(req request.PageRequest) (conveyances []models.Conveyance, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM conveyances"
	countQuery := "SELECT COUNT(*) FROM conveyances"
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
		var conveyance models.Conveyance
		err = rows.Scan(
			&conveyance.ID,
			&conveyance.Name,
			&conveyance.Code,
			&conveyance.Type,
			&conveyance.CreatedAt,
			&conveyance.CreatedByID,
			&conveyance.UpdatedAt,
			&conveyance.UpdatedByID,
			&conveyance.DeletedAt,
			&conveyance.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		conveyances = append(conveyances, conveyance)
	}

	return conveyances, total, err
}

func (repo *conveyanceRepository) GetAllConveyance() (conveyances []models.Conveyance, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			*
		FROM 
			conveyances
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conveyance models.Conveyance
		err = rows.Scan(
			&conveyance.ID,
			&conveyance.Name,
			&conveyance.Code,
			&conveyance.Type,
			&conveyance.CreatedAt,
			&conveyance.CreatedByID,
			&conveyance.UpdatedAt,
			&conveyance.UpdatedByID,
			&conveyance.DeletedAt,
			&conveyance.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		conveyances = append(conveyances, conveyance)
	}

	return conveyances, err
}

func (repo *conveyanceRepository) UpdateConveyance(id uuid.UUID, conveyanceReq dto.ToDBUpdateConveyance) (conveyanceRes *models.Conveyance, err error) {

	conveyanceRes = new(models.Conveyance)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE conveyances SET 
			name = $1, type = $2, updated_at = $3, updated_by = $4
		WHERE 
			id = $5 AND deleted_at IS NULL
		RETURNING 
			* `,
		conveyanceReq.Name,
		conveyanceReq.Type,
		updatedAtString,
		conveyanceReq.UpdatedByID,
		id,
	).Scan(
		&conveyanceRes.ID,
		&conveyanceRes.Name,
		&conveyanceRes.Code,
		&conveyanceRes.Type,
		&conveyanceRes.CreatedAt,
		&conveyanceRes.CreatedByID,
		&conveyanceRes.UpdatedAt,
		&conveyanceRes.UpdatedByID,
		&conveyanceRes.DeletedAt,
		&conveyanceRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conveyance with id %s not found", id)
		}

		return nil, err
	}

	return conveyanceRes, err
}

func (repo *conveyanceRepository) SoftDeleteConveyance(id uuid.UUID, conveyanceReq dto.ToDBDeleteConveyance) (conveyanceRes *models.Conveyance, err error) {

	conveyanceRes = new(models.Conveyance)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE conveyances SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			*`,
		deletedAtString,
		conveyanceReq.DeletedByID,
		id,
	).Scan(
		&conveyanceRes.ID,
		&conveyanceRes.Name,
		&conveyanceRes.Code,
		&conveyanceRes.Type,
		&conveyanceRes.CreatedAt,
		&conveyanceRes.CreatedByID,
		&conveyanceRes.UpdatedAt,
		&conveyanceRes.UpdatedByID,
		&conveyanceRes.DeletedAt,
		&conveyanceRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conveyance with id %s not found", id)
		}

		return nil, err
	}

	return conveyanceRes, err
}

func (repo *conveyanceRepository) CountConveyance() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			conveyances`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
