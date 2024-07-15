package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	category "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
)

type categoryRepository struct {
	Conn *sql.DB
}

func NewCategoryRepository(Conn *sql.DB) category.Repository {
	return &categoryRepository{Conn}
}

func (repo *categoryRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *categoryRepository) CreateCategory(trx *sql.Tx, categoryReq dto.ToDBCreateCategory) (categoryRes *models.Category, err error) {

	categoryRes = new(models.Category)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)
	queryString := `
	INSERT INTO categories
		(name, code, description, created_at, created_by)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		*`
	queryArgs := []interface{}{categoryReq.Name, categoryReq.Code, categoryReq.Description, createdAtString, categoryReq.CreatedByID}
	queryScanArgs := []interface{}{&categoryRes.ID, &categoryRes.Name, &categoryRes.Code, &categoryRes.Description, &categoryRes.CreatedAt, &categoryRes.CreatedByID, &categoryRes.UpdatedAt, &categoryRes.UpdatedByID, &categoryRes.DeletedAt, &categoryRes.DeletedByID}

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

	return categoryRes, err
}

func (repo *categoryRepository) GetCategoryByID(id uuid.UUID) (category *models.Category, err error) {
	category = new(models.Category)
	err = repo.Conn.QueryRow(
		`SELECT 
			*
		FROM 
			categories 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Code,
		&category.Description,
		&category.CreatedAt,
		&category.CreatedByID,
		&category.UpdatedAt,
		&category.UpdatedByID,
		&category.DeletedAt,
		&category.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with id %s not found", id)
		}

		return nil, err
	}

	return category, err
}

func (repo *categoryRepository) GetIndexCategory(req request.PageRequest) (categoryes []models.Category, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM categories"
	countQuery := "SELECT COUNT(*) FROM categories"
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
		var category models.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Code,
			&category.Description,
			&category.CreatedAt,
			&category.CreatedByID,
			&category.UpdatedAt,
			&category.UpdatedByID,
			&category.DeletedAt,
			&category.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		categoryes = append(categoryes, category)
	}

	return categoryes, total, err
}

func (repo *categoryRepository) GetAllCategory() (categoryes []models.Category, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			*
		FROM 
			categories
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.Code,
			&category.Description,
			&category.CreatedAt,
			&category.CreatedByID,
			&category.UpdatedAt,
			&category.UpdatedByID,
			&category.DeletedAt,
			&category.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		categoryes = append(categoryes, category)
	}

	return categoryes, err
}

func (repo *categoryRepository) UpdateCategory(id uuid.UUID, categoryReq dto.ToDBUpdateCategory) (categoryRes *models.Category, err error) {

	categoryRes = new(models.Category)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE categories SET 
			name = $1, description=$2, updated_at = $3, updated_by = $4
		WHERE 
			id = $5 AND deleted_at IS NULL
		RETURNING 
			*`,
		categoryReq.Name,
		categoryReq.Description,
		updatedAtString,
		categoryReq.UpdatedByID,
		id,
	).Scan(
		&categoryRes.ID,
		&categoryRes.Name,
		&categoryRes.Code,
		&categoryRes.Description,
		&categoryRes.CreatedAt,
		&categoryRes.CreatedByID,
		&categoryRes.UpdatedAt,
		&categoryRes.UpdatedByID,
		&categoryRes.DeletedAt,
		&categoryRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with id %s not found", id)
		}

		return nil, err
	}

	return categoryRes, err
}

func (repo *categoryRepository) SoftDeleteCategory(id uuid.UUID, categoryReq dto.ToDBDeleteCategory) (categoryRes *models.Category, err error) {

	categoryRes = new(models.Category)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE categories SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			*`,
		deletedAtString,
		categoryReq.DeletedByID,
		id,
	).Scan(
		&categoryRes.ID,
		&categoryRes.Name,
		&categoryRes.Code,
		&categoryRes.Description,
		&categoryRes.CreatedAt,
		&categoryRes.CreatedByID,
		&categoryRes.UpdatedAt,
		&categoryRes.UpdatedByID,
		&categoryRes.DeletedAt,
		&categoryRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category with id %s not found", id)
		}

		return nil, err
	}

	return categoryRes, err
}

func (repo *categoryRepository) CountCategory() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			categories`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
