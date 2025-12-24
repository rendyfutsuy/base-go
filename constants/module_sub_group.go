package constants

const (
	// SubGroup validation errors
	SubGroupNameAlreadyExists    = "Sub-group name already exists in this group"
	SubGroupCreateFailedIDNotSet = "failed to create sub-group: ID not set"
	SubGroupGroupNotFound        = "Goods group not found"
	SubGroupNotFound             = "sub-group with id %s not found"

	// Success messages
	SubGroupDeleteSuccess    = "Successfully deleted Sub Group"
	SubGroupStillUsedInTypes = "Sub-group is still used in active types"
)
