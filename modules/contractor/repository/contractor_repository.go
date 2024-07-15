package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	contractor "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type contractorRepository struct {
	Conn *sql.DB
}

func NewContractorRepository(Conn *sql.DB) contractor.Repository {
	return &contractorRepository{Conn}
}

func (repo *contractorRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *contractorRepository) CreateContractor(contractorReq dto.ToDBCreateContractor) (contractorRes *models.Contractor, err error) {

	contractorRes = new(models.Contractor)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO contractors
			(name, code, address, created_at, created_by)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING 
			* `,
		contractorReq.Name,
		contractorReq.Code,
		contractorReq.Address,
		createdAtString,
		contractorReq.CreatedByID,
	).Scan(
		&contractorRes.ID,
		&contractorRes.Name,
		&contractorRes.Code,
		&contractorRes.Address,
		&contractorRes.CreatedAt,
		&contractorRes.CreatedByID,
		&contractorRes.UpdatedAt,
		&contractorRes.UpdatedByID,
		&contractorRes.DeletedAt,
		&contractorRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return contractorRes, err
}

func (repo *contractorRepository) GetContractorByID(id uuid.UUID) (contractor *models.Contractor, err error) {
	contractor = new(models.Contractor)
	err = repo.Conn.QueryRow(
		`SELECT 
			*
		FROM 
			contractors 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&contractor.ID,
		&contractor.Name,
		&contractor.Code,
		&contractor.Address,
		&contractor.CreatedAt,
		&contractor.CreatedByID,
		&contractor.UpdatedAt,
		&contractor.UpdatedByID,
		&contractor.DeletedAt,
		&contractor.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contractor with id %s not found", id)
		}

		return nil, err
	}

	return contractor, err
}

func (repo *contractorRepository) GetIndexContractor(req request.PageRequest) (contractors []models.Contractor, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM contractors"
	countQuery := "SELECT COUNT(*) FROM contractors"
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
		var contractor models.Contractor
		err = rows.Scan(
			&contractor.ID,
			&contractor.Name,
			&contractor.Code,
			&contractor.Address,
			&contractor.CreatedAt,
			&contractor.CreatedByID,
			&contractor.UpdatedAt,
			&contractor.UpdatedByID,
			&contractor.DeletedAt,
			&contractor.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		contractors = append(contractors, contractor)
	}

	return contractors, total, err
}

func (repo *contractorRepository) GetAllContractor() (contractors []models.Contractor, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			*
		FROM 
			contractors
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var contractor models.Contractor
		err = rows.Scan(
			&contractor.ID,
			&contractor.Name,
			&contractor.Code,
			&contractor.Address,
			&contractor.CreatedAt,
			&contractor.CreatedByID,
			&contractor.UpdatedAt,
			&contractor.UpdatedByID,
			&contractor.DeletedAt,
			&contractor.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		contractors = append(contractors, contractor)
	}

	return contractors, err
}

func (repo *contractorRepository) UpdateContractor(id uuid.UUID, contractorReq dto.ToDBUpdateContractor) (contractorRes *models.Contractor, err error) {

	contractorRes = new(models.Contractor)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE contractors SET 
			name = $1, address = $2, updated_at = $3, updated_by = $4
		WHERE 
			id = $5 AND deleted_at IS NULL
		RETURNING 
			* `,
		contractorReq.Name,
		contractorReq.Address,
		updatedAtString,
		contractorReq.UpdatedByID,
		id,
	).Scan(
		&contractorRes.ID,
		&contractorRes.Name,
		&contractorRes.Code,
		&contractorRes.Address,
		&contractorRes.CreatedAt,
		&contractorRes.CreatedByID,
		&contractorRes.UpdatedAt,
		&contractorRes.UpdatedByID,
		&contractorRes.DeletedAt,
		&contractorRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contractor with id %s not found", id)
		}

		return nil, err
	}

	return contractorRes, err
}

func (repo *contractorRepository) SoftDeleteContractor(id uuid.UUID, contractorReq dto.ToDBDeleteContractor) (contractorRes *models.Contractor, err error) {

	contractorRes = new(models.Contractor)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE contractors SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			*`,
		deletedAtString,
		contractorReq.DeletedByID,
		id,
	).Scan(
		&contractorRes.ID,
		&contractorRes.Name,
		&contractorRes.Code,
		&contractorRes.Address,
		&contractorRes.CreatedAt,
		&contractorRes.CreatedByID,
		&contractorRes.UpdatedAt,
		&contractorRes.UpdatedByID,
		&contractorRes.DeletedAt,
		&contractorRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contractor with id %s not found", id)
		}

		return nil, err
	}

	return contractorRes, err
}

func (repo *contractorRepository) CountContractor() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			contractors`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
