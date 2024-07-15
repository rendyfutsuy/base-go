package utils

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/constants"
	"github.com/google/uuid"
)

// StringToUUID converts a string to a UUID.
func StringToUUID(request interface{}) (result uuid.UUID, err error) {
	// Try to assert the request as a string.
	stringUUID, ok := request.(string)
	if !ok {
		// If it's not a string, try to assert it as a UUID.
		uuid, ok := request.(uuid.UUID)
		if !ok {
			// If it's neither, return an error.
			err = constants.ErrorUUIDNotRecognized
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
