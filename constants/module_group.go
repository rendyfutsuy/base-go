package constants

const (
	// Group validation errors
	GroupNameAlreadyExists    = "Group name already exists"
	GroupCreateFailedIDNotSet = "failed to create group: ID not set"
	GroupStillUsedInSubGroups = "Group is still used in active sub-groups"
	GroupNotFound             = "group with id %s not found"

	// Success messages
	GroupDeleteSuccess = "Successfully deleted Group"
)
