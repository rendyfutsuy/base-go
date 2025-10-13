package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	"github.com/rendyfutsuy/base-go.git/models"
	account "github.com/rendyfutsuy/base-go.git/modules/account"
	"github.com/rendyfutsuy/base-go.git/modules/account/dto"
	"github.com/rendyfutsuy/base-go.git/utils"
)

type accountRepository struct {
	Conn *sql.DB
}

func NewAccountRepository(Conn *sql.DB) account.Repository {
	return &accountRepository{Conn}
}

func (repo *accountRepository) CreateTable(sqlFilePath string) (err error) {

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

func (repo *accountRepository) CreateAccount(accountReq dto.ToDBCreateAccount) (accountRes *models.Account, err error) {

	accountRes = new(models.Account)
	timeFormat := utils.ConfigVars.String("format.time")
	createdAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`INSERT INTO accounts
			(name, code, created_at, created_by)
		VALUES
			($1, $2, $3, $4)
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		accountReq.Name,
		accountReq.Code,
		createdAtString,
		accountReq.CreatedByID,
	).Scan(
		&accountRes.ID,
		&accountRes.Name,
		&accountRes.Code,
		&accountRes.CreatedAt,
		&accountRes.CreatedByID,
		&accountRes.UpdatedAt,
		&accountRes.UpdatedByID,
		&accountRes.DeletedAt,
		&accountRes.DeletedByID,
	)

	if err != nil {
		return nil, err
	}

	return accountRes, err
}

func (repo *accountRepository) GetAccountByID(id uuid.UUID) (account *models.Account, err error) {
	account = new(models.Account)
	err = repo.Conn.QueryRow(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			accounts 
		WHERE 
			id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&account.ID,
		&account.Name,
		&account.Code,
		&account.CreatedAt,
		&account.CreatedByID,
		&account.UpdatedAt,
		&account.UpdatedByID,
		&account.DeletedAt,
		&account.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with id %s not found", id)
		}

		return nil, err
	}

	return account, err
}

func (repo *accountRepository) GetIndexAccount(req request.PageRequest) (accounts []models.Account, total int, err error) {
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	baseQuery := "SELECT * FROM accounts"
	countQuery := "SELECT COUNT(*) FROM accounts"
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
		var account models.Account
		err = rows.Scan(
			&account.ID,
			&account.Name,
			&account.Code,
			&account.CreatedAt,
			&account.CreatedByID,
			&account.UpdatedAt,
			&account.UpdatedByID,
			&account.DeletedAt,
			&account.DeletedByID,
		)

		if err != nil {
			return nil, 0, err
		}

		accounts = append(accounts, account)
	}

	return accounts, total, err
}

func (repo *accountRepository) GetAllAccount() (accounts []models.Account, err error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
		FROM 
			accounts
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account
		err = rows.Scan(
			&account.ID,
			&account.Name,
			&account.Code,
			&account.CreatedAt,
			&account.CreatedByID,
			&account.UpdatedAt,
			&account.UpdatedByID,
			&account.DeletedAt,
			&account.DeletedByID,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, err
}

func (repo *accountRepository) UpdateAccount(id uuid.UUID, accountReq dto.ToDBUpdateAccount) (accountRes *models.Account, err error) {

	accountRes = new(models.Account)
	timeFormat := utils.ConfigVars.String("format.time")
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE accounts SET 
			name = $1, updated_at = $2, updated_by = $3
		WHERE 
			id = $4 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		accountReq.Name,
		updatedAtString,
		accountReq.UpdatedByID,
		id,
	).Scan(
		&accountRes.ID,
		&accountRes.Name,
		&accountRes.Code,
		&accountRes.CreatedAt,
		&accountRes.CreatedByID,
		&accountRes.UpdatedAt,
		&accountRes.UpdatedByID,
		&accountRes.DeletedAt,
		&accountRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with id %s not found", id)
		}

		return nil, err
	}

	return accountRes, err
}

func (repo *accountRepository) SoftDeleteAccount(id uuid.UUID, accountReq dto.ToDBDeleteAccount) (accountRes *models.Account, err error) {

	accountRes = new(models.Account)
	timeFormat := utils.ConfigVars.String("format.time")
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE accounts SET 
			deleted_at = $1, deleted_by = $2
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id, name, code, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by`,
		deletedAtString,
		accountReq.DeletedByID,
		id,
	).Scan(
		&accountRes.ID,
		&accountRes.Name,
		&accountRes.Code,
		&accountRes.CreatedAt,
		&accountRes.CreatedByID,
		&accountRes.UpdatedAt,
		&accountRes.UpdatedByID,
		&accountRes.DeletedAt,
		&accountRes.DeletedByID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with id %s not found", id)
		}

		return nil, err
	}

	return accountRes, err
}

func (repo *accountRepository) CountAccount() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			accounts`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}
