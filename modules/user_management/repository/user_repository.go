package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

// CreateUser creates a new user information entry in the database.
//
// It takes a ToDBCreateUser parameter and returns an User pointer and an error.
func (repo *userRepository) CreateUser(userReq dto.ToDBCreateUser) (userRes *models.User, err error) {

	// initialize: user user model, time format to created at string,
	userRes = new(models.User)
	timeFormat := constants.FormatTimezone
	createdAtString := time.Now().UTC().Format(timeFormat)
	ExpiredAtString := time.Now().UTC().AddDate(0, 3, 0).Format(timeFormat)

	// execute query to insert user user
	// assign return value to userRes variable
	err = repo.Conn.QueryRow(
		`INSERT INTO users
			(full_name, email, role_id, is_active, gender, api_key, created_at, updated_at, password, password_expired_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING 
			id, full_name, created_at, updated_at, deleted_at`,
		userReq.FullName,
		userReq.Email,
		userReq.RoleId,
		userReq.IsActive,
		userReq.Gender,
		userReq.ApiKey,
		createdAtString,
		createdAtString,
		"temp",
		ExpiredAtString,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)
	// if error occurs, return error
	if err != nil {
		return nil, err
	}

	return userRes, err
}

// GetUserByID retrieves an user information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an User pointer and an error.
func (repo *userRepository) GetUserByID(id uuid.UUID) (user *models.User, err error) {
	// initialize user variable
	user = new(models.User)

	// fetch data from database by id that passed
	// assign return value to user variable
	err = repo.Conn.QueryRow(
		`SELECT 
			usr.id,
			usr.full_name,
			usr.email,
			usr.created_at,
			usr.updated_at,
			usr.deleted_at,
			usr.role_id,
			usr.is_active,
			rl.name,
			usr.api_key,
			usr.gender,
			CASE 
				WHEN usr.is_active THEN 'active'
				ELSE 'inactive'
			END AS active_status,
			CASE 
				WHEN usr.counter >= 3 THEN true
				ELSE false
			END AS is_blocked
		FROM 
			users usr
		JOIN
			roles rl
		ON
			usr.role_id = rl.id
		WHERE 
			usr.id = $1 AND usr.deleted_at IS NULL`,
		id,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.RoleId,
		&user.IsActive,
		&user.RoleName,
		&user.ApiKey,
		&user.Gender,
		&user.ActiveStatus,
		&user.IsBlocked,
	)

	return user, err
}

// GetIndexUser retrieves a paginated list of user information from the database.
//
// It takes a PageRequest parameter and returns a slice of User, the total number of
// user information entries, and an error.
// its can search by user name, user code, user alias_1, user alias_2, user alias_3, user alias_4, user address, user email, user phone_number, type name
func (repo *userRepository) GetIndexUser(req request.PageRequest, filter dto.ReqUserIndexFilter) (users []models.User, total int, err error) {
	// initialize: pagination page, search query to local variable
	offSet := (req.Page - 1) * req.PerPage
	searchQuery := req.Search

	// Construct the SQL query
	JoinFill := ` JOIN
		roles rl
		ON
			usr.role_id = rl.id
		`

	// construct select fillable
	selectFill := `
		usr.id,
		usr.full_name,
		usr.email,
		usr.gender,
		usr.is_active,
		usr.counter,
		usr.created_at,
		usr.updated_at,
		usr.deleted_at,
		CASE 
			WHEN usr.is_active THEN 'active'
			ELSE 'inactive'
		END AS active_status,
		CASE 
			WHEN usr.counter >= 3 THEN true
			ELSE false
		END AS is_blocked,
		rl.name AS role_name
	`

	// append select and join query to base query
	baseQuery := "SELECT " + selectFill + " FROM users usr" + JoinFill

	// append join query to count query
	countQuery := "SELECT COUNT(DISTINCT usr.id) FROM users usr" + JoinFill

	// initialize common query for deleted condition
	whereClause := " WHERE usr.deleted_at IS NULL"

	// assign search query, based on searchable field.
	if searchQuery != "" {
		// can search by full_name, gender, email and role_name
		searchUser := "usr.full_name ILIKE '%' || $1 || '%' OR usr.gender ILIKE '%' || $1 || '%' OR usr.email ILIKE '%' || $1 || '%' OR rl.name ILIKE '%' || $1 || '%'"
		whereClause += " AND (" + searchUser + ")"
	}

	params := []interface{}{}
	paramIndex := 1

	if searchQuery != "" {
		paramIndex = 2
	}

	// multiple input filter - BEGIN
	if len(filter.RoleIds) > 0 {
		whereClause += fmt.Sprintf(" AND rl.id = ANY($%d)", paramIndex)
		params = append(params, pq.Array(filter.RoleIds))
		paramIndex++
	}
	// multiple input filter - END

	// Single input filter - BEGIN
	if filter.RoleName != "" {
		whereClause += fmt.Sprintf(" AND rl.name = $%d", paramIndex)
		params = append(params, filter.RoleName) // Adding % for partial match
		paramIndex++
	}
	// Single input filter - END

	// Initialize Default sorting
	sortBy := "usr.created_at"
	sortOrder := "DESC" // Sort from newest to oldest
	if req.SortBy != "" {
		// intercept sort by
		sortBy = repo.SortColumnMapping(req.SortBy)

		// adjust sort order
		if req.SortOrder != "" {
			sortOrder = req.SortOrder
		}
	}

	// initialize Default Sorting Query
	orderClause := " ORDER BY " + sortBy + " " + sortOrder
	limitClause := fmt.Sprintf(" LIMIT %d OFFSET %d", req.PerPage, offSet)

	// count total
	if searchQuery != "" {
		err = repo.Conn.QueryRow(countQuery+whereClause, append([]interface{}{searchQuery}, params...)...).Scan(&total)
	} else {
		err = repo.Conn.QueryRow(countQuery+whereClause, params...).Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}

	// retrieve paginated
	rows := new(sql.Rows)
	if searchQuery != "" {
		rows, err = repo.Conn.Query(baseQuery+whereClause+orderClause+limitClause, append([]interface{}{searchQuery}, params...)...)
	} else {
		rows, err = repo.Conn.Query(baseQuery+whereClause+orderClause+limitClause, params...)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// assign pagination as models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Gender,
			&user.IsActive,
			&user.Counter,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.ActiveStatus,
			&user.IsBlocked,
			&user.RoleName,
		)

		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	return users, total, err
}

// GetAllUser retrieves all user information entries from the database.
//
// Returns a slice of models.User and an error.
func (repo *userRepository) GetAllUser() ([]models.User, error) {
	rows, err := repo.Conn.Query(
		`SELECT 
			id,
			full_name,
			created_at
		FROM 
			users
		WHERE
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.FullName,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user information entry in the database.
//
// It takes an ID of the user information and a ToDBUpdateUser parameter.
// It returns an User pointer and an error.
//
// The function updates the user information in the database with the provided ID.
// It sets the name, updated_at, updated_by, email, phone_number, user_type_id,
// alias_1, alias_2, alias_3, alias_4, and address fields of the user information.
// If the user information with the provided ID is not found, it returns an error.
// If there is an error during the update, it returns the error.
func (repo *userRepository) UpdateUser(id uuid.UUID, userReq dto.ToDBUpdateUser) (userRes *models.User, err error) {
	userRes = new(models.User)
	timeFormat := constants.FormatTimezone
	updatedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			full_name = $1,
			updated_at = $2,
			email = $4,
			api_key = $5,
			gender = $6,
			is_active = $7,
			role_id = $8
		WHERE 
			id = $3 AND deleted_at IS NULL
		RETURNING 
			id,
			full_name,
			created_at,
			updated_at,
			deleted_at`,
		userReq.FullName,
		updatedAtString,
		id,
		userReq.Email,
		userReq.ApiKey,
		userReq.Gender,
		userReq.IsActive,
		userReq.RoleId,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

// SoftDeleteUser soft deletes an user user entry in the database.
//
// It takes an id of type uuid.UUID and an userReq of type dto.ToDBDeleteUser as parameters.
// It returns the soft deleted user user entry of type models.User and an error.
func (repo *userRepository) SoftDeleteUser(id uuid.UUID, userReq dto.ToDBDeleteUser) (userRes *models.User, err error) {

	userRes = new(models.User)
	timeFormat := constants.FormatTimezone
	deletedAtString := time.Now().UTC().Format(timeFormat)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			deleted_at = $1
		WHERE 
			id = $2 AND deleted_at IS NULL
		RETURNING 
			id, name, created_at, updated_at, deleted_at`,
		deletedAtString,
		id,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

func (repo *userRepository) BlockUser(id uuid.UUID) (userRes *models.User, err error) {

	userRes = new(models.User)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			counter = 4
		WHERE 
			id = $1
		RETURNING 
			id, full_name, counter, created_at, updated_at, deleted_at`,
		id,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.Counter,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

func (repo *userRepository) UnBlockUser(id uuid.UUID) (userRes *models.User, err error) {

	userRes = new(models.User)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			counter = 0
		WHERE 
			id = $1
		RETURNING 
			id, full_name, counter, created_at, updated_at, deleted_at`,
		id,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.Counter,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

func (repo *userRepository) ActivateUser(id uuid.UUID) (userRes *models.User, err error) {

	userRes = new(models.User)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			is_active = true
		WHERE 
			id = $1
		RETURNING 
			id, full_name, is_active, created_at, updated_at, deleted_at`,
		id,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.IsActive,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

func (repo *userRepository) DisActivateUser(id uuid.UUID) (userRes *models.User, err error) {

	userRes = new(models.User)

	err = repo.Conn.QueryRow(
		`UPDATE users SET 
			is_active = false
		WHERE 
			id = $1
		RETURNING 
			id, full_name, is_active, created_at, updated_at, deleted_at`,
		id,
	).Scan(
		&userRes.ID,
		&userRes.FullName,
		&userRes.IsActive,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
		&userRes.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}

		return nil, err
	}

	return userRes, err
}

// CountUser retrieves the count of user information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *userRepository) CountUser() (count *int, err error) {
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*)
		FROM 
			users`,
	).Scan(&count)

	if err != nil {
		return nil, err
	}

	return count, err
}

// EmailIsNotDuplicated checks if an email is not duplicated in the users table, excluding a specific ID if provided.
//
// Parameters:
// - email: the email to check for duplication.
// - excludedId: the ID to exclude from the check. If set to uuid.Nil, no exclusion is applied.
//
// Returns:
// - bool: true if the email is not duplicated, false otherwise.
// - error: an error if the check fails.
func (repo *userRepository) EmailIsNotDuplicated(email string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			users
		WHERE
			email = $1 AND deleted_at IS NULL`

	params := []interface{}{email}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	result := 0

	// assert email is nt duplicated
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

// ApiKeyIsNotDuplicated checks if the provided API key is not duplicated in the users table.
//
// Parameters:
// - apiKey: the API key to check for duplication.
// - excludedId: the ID to exclude from the check. If set to uuid.Nil, no exclusion is applied.
//
// Returns:
// - bool: true if the API key is not duplicated, false otherwise.
// - error: an error if the check fails.
func (repo *userRepository) ApiKeyIsNotDuplicated(apiKey utils.NullString, excludedId uuid.UUID) (bool, error) {

	// if input not valid return true without error
	if apiKey.String == "" {
		return true, nil
	}

	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			users
		WHERE
			api_key = $1 AND deleted_at IS NULL`

	params := []interface{}{apiKey}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	result := 0

	// assert apiKey is nt duplicated
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

// UserNameIsNotDuplicated checks if the provided user name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *userRepository) UserNameIsNotDuplicated(name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			users
		WHERE
			full_name = $1 AND deleted_at IS NULL`

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

// GetDuplicatedUser retrieves the user information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the user information to retrieve.
// - excludedId: the ID of the user information to exclude from the result.
//
// Returns:
// - user: a pointer to the retrieved user information.
// - err: an error if there was a problem retrieving the user information.
func (repo *userRepository) GetDuplicatedUser(name string, excludedId uuid.UUID) (user *models.User, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			users
		WHERE
			full_name = $1 AND deleted_at IS NULL`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize user
	user = &models.User{}

	// assert name is not duplicated
	err = repo.Conn.QueryRow(baseQuery, params...).Scan(
		&user.ID,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated User
	return user, nil
}

// UserNameIsNotDuplicatedOnSoftDeleted checks if the provided user name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *userRepository) UserNameIsNotDuplicatedOnSoftDeleted(name string, excludedId uuid.UUID) (bool, error) {
	baseQuery := `SELECT 
			COUNT(*)
		FROM 
			users
		WHERE
			full_name = $1`

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

// GetDuplicatedUserOnSoftDeleted retrieves the user information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the user information to retrieve.
// - excludedId: the ID of the user information to exclude from the result.
//
// Returns:
// - user: a pointer to the retrieved user information.
// - err: an error if there was a problem retrieving the user information.
func (repo *userRepository) GetDuplicatedUserOnSoftDeleted(name string, excludedId uuid.UUID) (user *models.User, err error) {
	baseQuery := `SELECT 
			id, name, created_at, updated_at
		FROM 
			users
		WHERE
			full_name = $1`

	params := []interface{}{name}

	if excludedId != uuid.Nil {
		baseQuery += ` AND id <> $2`
		params = append(params, excludedId)
	}

	// Initialize user
	user = &models.User{}

	// assert name is not duplicated
	err = repo.Conn.QueryRow(baseQuery, params...).Scan(
		&user.ID,
		&user.FullName,
		&user.CreatedAt,
		&user.UpdatedAt)

	// if have error, return nil and error
	if err != nil {
		return nil, err
	}

	// return Duplicated User
	return user, nil
}

func (repo *userRepository) SortColumnMapping(selectedSortLabel string) string {
	response := ""
	sortLabels := map[string][]string{
		"id": []string{
			"id",
		},
		"full_name": []string{
			"full_name",
			"name",
		},
		"email": []string{
			"email",
		},
		"gender": []string{
			"gender",
		},
		"is_active": []string{
			"is_active",
		},
		"counter": []string{
			"counter",
		},
		"created_at": []string{
			"created_at",
		},
		"updated_at": []string{
			"updated_at",
		},
		"deleted_at": []string{
			"deleted_at",
		},
		"active_status": []string{
			"active_status",
		},
		"is_blocked": []string{
			"is_blocked",
		},
		"role_name": []string{
			"role_name",
		},
	}

	// Loop through the map
	for DBcolumn, sortLabels := range sortLabels {
		for _, sortLabel := range sortLabels {
			// Check if the current sort label matches the selected sort label
			if sortLabel == selectedSortLabel {
				response = DBcolumn
			}
		}
	}

	return response
}
