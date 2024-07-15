package repository

import (
	"database/sql"
	"fmt"
	"math"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// shipyardTable is the name of the table in the database that stores shipyard data.
	shipyardTable = "shipyards"
)

// shipyardRepository is the struct that implements the Repository interface for shipyard.
type shipyardRepository struct {
	Conn *sql.DB // Conn is the database connection object.
}

// NewShipyardRepository returns a new instance of a shipyardRepository.
func NewShipyardRepository(Conn *sql.DB) shipyard.Repository {
	return &shipyardRepository{Conn}
}

// CreateShipyard inserts a new shipyard record in the database.
func (repo *shipyardRepository) CreateShipyard(data *models.Shipyard) (err error) {
	// The SQL query string for inserting a new shipyard record.
	// It uses placeholders (%s and $1, $2, $3) for the table name and column values.
	query := `
		INSERT INTO %s
			(name, code, yard, created_at, updated_at)
		VALUES 	
			($1, $2, $3, NOW(), NOW())
		RETURNING id, code,	created_at, updated_at
	`
	// The parameters for the SQL query.
	param := []any{data.Name, data.Code, data.Yard}

	// Execute the SQL query with the parameters and scan the result into the data object.
	err = repo.Conn.QueryRow(fmt.Sprintf(query, shipyardTable), param...).Scan(&data.ID, &data.Code, &data.CreatedAt, &data.UpdatedAt)

	return
}

// UpdateShipyard updates an existing shipyard record in the database.
func (repo *shipyardRepository) UpdateShipyard(data *models.Shipyard) (err error) {
	// Define the SQL query with placeholders for the table name and column values.
	query := `
		UPDATE %s
		SET 
			name = $1, 
			yard = $2,
			updated_at = NOW()
		WHERE 
			id = $3
		RETURNING code,updated_at	
	`
	// Check if the ID of the shipyard is empty.
	if data.ID == uuid.Nil {
		// If the ID is empty, log an error and return.
		err = errorUUIDIsEmpty
		zap.S().Error(err)
		return
	}

	// Define the parameters for the SQL query.
	param := []any{data.Name, data.Yard, data.ID}

	// Execute the SQL query with the parameters and scan the result into the data object.
	err = repo.Conn.QueryRow(fmt.Sprintf(query, shipyardTable), param...).Scan(&data.Code, &data.UpdatedAt)

	return
}

// DeleteShipyard deletes a shipyard record from the database.
func (repo *shipyardRepository) DeleteShipyard(data *models.Shipyard) (err error) {
	// Define the SQL query with placeholders for the table name and column values.
	query := `
		DELETE FROM %s
		WHERE 
			id = $1	
	`
	// Check if the ID of the shipyard is empty.
	if data.ID == uuid.Nil {
		// If the ID is empty, log an error and return.
		err = errorUUIDIsEmpty
		zap.S().Error(err)
		return
	}

	// Define the parameters for the SQL query.
	param := []any{data.ID}

	// Execute the SQL query with the parameters.
	// If an error occurs during the execution of the SQL query, it is returned.
	err = repo.Conn.QueryRow(fmt.Sprintf(query, shipyardTable), param...).Err()

	return
}

// StringToUUID converts a string to a UUID.
func StringToUUID(request interface{}) (result uuid.UUID, err error) {
	// Try to assert the request as a string.
	stringUUID, ok := request.(string)
	if !ok {
		// If it's not a string, try to assert it as a UUID.
		uuid, ok := request.(uuid.UUID)
		if !ok {
			// If it's neither, return an error.
			err = errorUUIDNotRecognized
			return
		}
		// If it's a UUID, return it.
		result = uuid
	} else {
		// If it's a string, try to parse it as a UUID.
		uuid, parseErr := uuid.Parse(stringUUID)
		if parseErr != nil {
			// If the parsing fails, return an error.
			err = fmt.Errorf("requested param is string")
			return
		}
		// If the parsing succeeds, return the parsed UUID.
		result = uuid
	}
	return
}

// FindShipyardByUUIDOrCode finds a shipyard record by its UUID or code.
func (repo *shipyardRepository) FindShipyardByUUIDOrCode(request interface{}) (result models.Shipyard, err error) {
	// Define the SQL query with placeholders for the table name and column values.
	query := `
		SELECT
			id, name, code, yard, created_at, updated_at
		FROM
			%s
		WHERE
			%s
		AND	deleted_at IS NULL
	`

	whereQuery := ""
	param := []any{}
	scanField := []any{&result.ID, &result.Name, &result.Code, &result.Yard, &result.CreatedAt, &result.UpdatedAt}

	// Try to convert the request to a UUID.
	uuid, parseErr := StringToUUID(request)
	if parseErr != errorUUIDNotRecognized {
		if parseErr == nil {
			// If the request is a valid UUID, use it in the WHERE clause.
			whereQuery = "(id = $1)"
			param = append(param, uuid)
		} else {
			// If the request is not a valid UUID, use it as a code in the WHERE clause.
			whereQuery = "(code = $1)"
			param = append(param, request)
		}
	} else {
		// If the request is neither a UUID nor a code, log an error and return.
		err = parseErr
		zap.S().Error(err)
		return
	}
	// Replace the placeholders in the query with the table name and WHERE clause.
	query = fmt.Sprintf(query, shipyardTable, whereQuery)
	// Execute the SQL query and scan the result into the result object.
	err = repo.Conn.QueryRow(query, param...).Scan(scanField...)

	// Handle any errors that occurred during the execution of the SQL query.
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were returned, log an error and return.
			err = fmt.Errorf(errorNotFound, err)
			zap.S().Error(err)
			return
		} else {
			// If another error occurred, log it and return.
			err = fmt.Errorf(errorQueryScan, err)
			zap.S().Error(err)
			return
		}
	}

	return
}

// FetchShipyards fetches all shipyard records that match the provided query request.
func (repo *shipyardRepository) FetchShipyards(request shipyard.QueryRequest) (result []*models.Shipyard, total int, lastPage int, err error) {
	// Define the SQL query with placeholders for the table name and column values.
	query := `
		SELECT 
			id, name, code, yard, created_at, updated_at
		FROM 
			%s
		WHERE 
			%s
	`

	// Add sorting if the sort field is provided
	sortField, sortOrder := request.GetSort()
	if sortField != "" {
		query += `
		ORDER BY ` + sortField + ` ` + sortOrder
	}

	var args []any
	// Get the parameters for the WHERE clause from the request.
	args = request.GetParam()

	paginate, limit, offset := request.IsPaginate()
	if paginate {
		lengthParam := len(args)
		query += fmt.Sprintf(`
		LIMIT $%d OFFSET $%d
		`, lengthParam+1, lengthParam+2)
		args = append(args, limit, offset)
	}

	var rows *sql.Rows

	// Execute the SQL query and get the result rows.
	if len(args) > 0 {
		rows, err = repo.Conn.Query(fmt.Sprintf(query, shipyardTable, request.GetCondition()), args...)
	} else {
		rows, err = repo.Conn.Query(fmt.Sprintf(query, shipyardTable, request.GetCondition()))
	}
	if err != nil {
		// If an error occurred, log it and return.
		zap.S().Error(err)
		return
	}
	// Ensure the rows are closed after we're done with them.
	defer rows.Close()

	// Iterate over the rows.
	for rows.Next() {
		var s models.Shipyard
		// Scan the values from the row into the Shipyard object.
		err = rows.Scan(&s.ID, &s.Name, &s.Code, &s.Yard, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			// If an error occurred, log it and return.
			zap.S().Error(err)
			return
		}
		// Append the Shipyard object to the result slice.
		result = append(result, &s)
	}

	// Check for any error that occurred while iterating over the rows.
	if err = rows.Err(); err != nil {
		zap.S().Error(err)
		return
	}

	// If pagination is enabled, it executes another SQL query to get the total count of rows.
	if paginate {
		countQuery := `
			SELECT 
				COUNT(*)
			FROM 
				%s
			WHERE 
				%s
		`
		var totalRow *sql.Row

		if len(request.GetParam()) > 0 {
			totalRow = repo.Conn.QueryRow(fmt.Sprintf(countQuery, shipyardTable, request.GetCondition()), request.GetParam()...)
		} else {
			totalRow = repo.Conn.QueryRow(fmt.Sprintf(countQuery, shipyardTable, request.GetCondition()))
		}
		err = totalRow.Scan(&total)
		if err != nil {
			zap.S().Error(err)
			return
		}

		// Calculate the last page
		lastPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return
}
