package repository

import "fmt"

const (
	errorNotFound  = "Requested shipyard not found : %w" // Error message when a requested shipyard is not found.
	errorQueryScan = "QueryRow scan error : %w"          // Error message when there's an error scanning the result of a QueryRow operation.
)

var (
	errorUUIDNotRecognized = fmt.Errorf("Requested param is not recognized") // Error when a requested parameter is not recognized.
	errorUUIDIsEmpty       = fmt.Errorf("Requested shipyard is not exists")  // Error when a requested shipyard does not exist.
)
