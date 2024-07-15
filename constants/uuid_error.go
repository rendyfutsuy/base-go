package constants

import "fmt"

var (
	ErrorUUIDNotRecognized = fmt.Errorf("Requested param is not recognized") // Error when a requested parameter is not recognized.
	ErrorUUIDIsEmpty       = fmt.Errorf("Requested shipyard is not exists")  // Error when a requested shipyard does not exist.
)
