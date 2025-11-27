package constants

const (
	// Type validation errors
	TypeNameAlreadyExists    = "Type name already exists in this subgroup"
	TypeCreateFailedIDNotSet = "failed to create type: ID not set"
	TypeSubGroupNotFound     = "Sub-group not found"
	TypeNotFound             = "type with id %s not found"

	// Success messages
	TypeDeleteSuccess       = "Successfully deleted Type"
	TypeStillUsedInBackings = "Type is still used in active backings"
)
