package repository

import (
	"database/sql"
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	cobsubcob "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/dto"
	"github.com/google/uuid"
)

type cobRepository struct {
	Conn *sql.DB
}

func NewCobSubcobRepository(conn *sql.DB) cobsubcob.Repository {
	return &cobRepository{conn}
}

func (repo *cobRepository) StartTransaction() (*sql.Tx, error) {
	return repo.Conn.Begin()
}

func (repo *cobRepository) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (repo *cobRepository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (repo *cobRepository) CreateCob(trx *sql.Tx, cobReq dto.ToDBCreateCob) (cobRes *models.Cob, err error) {
	cobRes = new(models.Cob)

	queryString := `
	INSERT INTO cobs
		(category_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_by)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING
		*`
	queryArgs := []interface{}{
		cobReq.CategoryID,
		cobReq.Name,
		cobReq.Code,
		cobReq.Forms,
		cobReq.ActiveDate,
		cobReq.IsHiddenFromFacultative,
		cobReq.IsInactive,
		cobReq.IsFromWebCrediit,
		cobReq.CreatedByID,
	}
	queryScanArgs := []interface{}{
		&cobRes.ID,
		&cobRes.CategoryID,
		&cobRes.Name,
		&cobRes.Code,
		&cobRes.Forms,
		&cobRes.ActiveDate,
		&cobRes.IsHiddenFromFacultative,
		&cobRes.IsInactive,
		&cobRes.IsFromWebCrediit,
		&cobRes.CreatedAt,
		&cobRes.CreatedByID,
		&cobRes.UpdatedAt,
		&cobRes.UpdatedByID,
		&cobRes.DeletedAt,
		&cobRes.DeletedByID,
	}

	if trx == nil {

		err = repo.Conn.QueryRow(
			queryString,
			queryArgs...,
		).Scan(
			queryScanArgs...,
		)
	} else {
		err = trx.QueryRow(
			queryString,
			queryArgs...,
		).Scan(
			queryScanArgs...,
		)
	}

	if err != nil {
		return nil, err
	}

	return cobRes, err
}

func (repo *cobRepository) CreateSubcob(trx *sql.Tx, subcobReq dto.ToDBCreateSubcob) (subcobRes *models.Subcob, err error) {
	subcobRes = new(models.Subcob)

	queryString := `
	INSERT INTO subcobs
		(category_id, cob_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_by)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING
		*`
	queryArgs := []interface{}{
		subcobReq.CategoryID,
		subcobReq.CobID,
		subcobReq.Name,
		subcobReq.Code,
		subcobReq.Forms,
		subcobReq.ActiveDate,
		subcobReq.IsHiddenFromFacultative,
		subcobReq.IsInactive,
		subcobReq.IsFromWebCrediit,
		subcobReq.CreatedByID,
	}
	queryScanArgs := []interface{}{
		&subcobRes.ID,
		&subcobRes.CategoryID,
		&subcobRes.CobID,
		&subcobRes.Name,
		&subcobRes.Code,
		&subcobRes.Forms,
		&subcobRes.ActiveDate,
		&subcobRes.IsHiddenFromFacultative,
		&subcobRes.IsInactive,
		&subcobRes.IsFromWebCrediit,
		&subcobRes.CreatedAt,
		&subcobRes.CreatedByID,
		&subcobRes.UpdatedAt,
		&subcobRes.UpdatedByID,
		&subcobRes.DeletedAt,
		&subcobRes.DeletedByID,
	}

	if trx == nil {
		err = repo.Conn.QueryRow(
			queryString,
			queryArgs...,
		).Scan(
			queryScanArgs...,
		)
	} else {
		err = trx.QueryRow(
			queryString,
			queryArgs...,
		).Scan(
			queryScanArgs...,
		)
	}

	if err != nil {
		return nil, err
	}

	return subcobRes, nil
}

func (repo *cobRepository) GetCobByID(id uuid.UUID) (cobRes *models.Cob, err error) {
	cobRes = new(models.Cob)

	err = repo.Conn.QueryRow(
		`SELECT
			id, category_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM
			cobs
		WHERE
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&cobRes.ID,
		&cobRes.CategoryID,
		&cobRes.Name,
		&cobRes.Code,
		&cobRes.Forms,
		&cobRes.ActiveDate,
		&cobRes.IsHiddenFromFacultative,
		&cobRes.IsInactive,
		&cobRes.IsFromWebCrediit,
		&cobRes.CreatedAt,
		&cobRes.CreatedByID,
		&cobRes.UpdatedAt,
		&cobRes.UpdatedByID,
		&cobRes.DeletedAt,
		&cobRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return cobRes, err
}

func (repo *cobRepository) GetSubcobByID(id uuid.UUID) (subcobRes *models.Subcob, err error) {
	subcobRes = new(models.Subcob)

	err = repo.Conn.QueryRow(
		`SELECT
			id, category_id, cob_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM
			subcobs
		WHERE
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&subcobRes.ID,
		&subcobRes.CategoryID,
		&subcobRes.CobID,
		&subcobRes.Name,
		&subcobRes.Code,
		&subcobRes.Forms,
		&subcobRes.ActiveDate,
		&subcobRes.IsHiddenFromFacultative,
		&subcobRes.IsInactive,
		&subcobRes.IsFromWebCrediit,
		&subcobRes.CreatedAt,
		&subcobRes.CreatedByID,
		&subcobRes.UpdatedAt,
		&subcobRes.UpdatedByID,
		&subcobRes.DeletedAt,
		&subcobRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return subcobRes, err
}

func (repo *cobRepository) GetIndexCob(req request.PageRequest) (cobs []models.Cob, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search
	filter := req.Filters

	// Construct the SQL query
	baseQuery := "SELECT * FROM cobs"
	countQuery := "SELECT COUNT(*) FROM cobs"
	whereClause := " WHERE deleted_at IS NULL"

	// set filter
	if len(filter) > 0 {
		whereClause += " AND " + fromFilterToWhere(filter)
	}

	// set search
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
		var cob models.Cob
		err = rows.Scan(
			&cob.ID,
			&cob.CategoryID,
			&cob.Name,
			&cob.Code,
			&cob.Forms,
			&cob.ActiveDate,
			&cob.IsHiddenFromFacultative,
			&cob.IsInactive,
			&cob.IsFromWebCrediit,
			&cob.CreatedAt,
			&cob.CreatedByID,
			&cob.UpdatedAt,
			&cob.UpdatedByID,
			&cob.DeletedAt,
			&cob.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		cobs = append(cobs, cob)
	}

	return cobs, total, err
}

func (repo *cobRepository) GetIndexSubcob(req request.PageRequest) (Subcobs []models.Subcob, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search
	filter := req.Filters

	// Construct the SQL query
	baseQuery := "SELECT * FROM subcobs"
	countQuery := "SELECT COUNT(*) FROM subcobs"
	whereClause := " WHERE deleted_at IS NULL"

	// set filter
	if len(filter) > 0 {
		whereClause += " AND " + fromFilterToWhere(filter)
	}

	// set search
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
		var subcob models.Subcob
		err = rows.Scan(
			&subcob.ID,
			&subcob.CategoryID,
			&subcob.CobID,
			&subcob.Name,
			&subcob.Code,
			&subcob.Forms,
			&subcob.ActiveDate,
			&subcob.IsHiddenFromFacultative,
			&subcob.IsInactive,
			&subcob.IsFromWebCrediit,
			&subcob.CreatedAt,
			&subcob.CreatedByID,
			&subcob.UpdatedAt,
			&subcob.UpdatedByID,
			&subcob.DeletedAt,
			&subcob.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		Subcobs = append(Subcobs, subcob)
	}

	return Subcobs, total, err
}

func (repo *cobRepository) GetAllCob() (cobs []models.Cob, err error) {
	rows, err := repo.Conn.Query(
		`SELECT
			id, category_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM
			cobs
		WHERE
			deleted_at IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cob models.Cob
		err = rows.Scan(
			&cob.ID,
			&cob.CategoryID,
			&cob.Name,
			&cob.Code,
			&cob.Forms,
			&cob.ActiveDate,
			&cob.IsHiddenFromFacultative,
			&cob.IsInactive,
			&cob.IsFromWebCrediit,
			&cob.CreatedAt,
			&cob.CreatedByID,
			&cob.UpdatedAt,
			&cob.UpdatedByID,
			&cob.DeletedAt,
			&cob.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		cobs = append(cobs, cob)
	}

	return cobs, err
}

func (repo *cobRepository) GetAllSubcob() (subcobs []models.Subcob, err error) {
	rows, err := repo.Conn.Query(
		`SELECT
			id, category_id, cob_id, name, code, forms, active_date, is_hidden_from_facultative, is_inactive, is_from_web_credit, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM
			subcobs
		WHERE
			deleted_at IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subcob models.Subcob
		err = rows.Scan(
			&subcob.ID,
			&subcob.CategoryID,
			&subcob.CobID,
			&subcob.Name,
			&subcob.Code,
			&subcob.Forms,
			&subcob.ActiveDate,
			&subcob.IsHiddenFromFacultative,
			&subcob.IsInactive,
			&subcob.IsFromWebCrediit,
			&subcob.CreatedAt,
			&subcob.CreatedByID,
			&subcob.UpdatedAt,
			&subcob.UpdatedByID,
			&subcob.DeletedAt,
			&subcob.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		subcobs = append(subcobs, subcob)
	}

	return subcobs, err
}
