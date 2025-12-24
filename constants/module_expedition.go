package constants

import "fmt"

const (
	ExpeditionContactTypeTelp  = "telp"
	ExpeditionContactTypePhone = "hp"
)

const (
	// Expedition validation errors
	ExpeditionNameAlreadyExists    = "Expedition name already exists"
	ExpeditionCreateFailedIDNotSet = "failed to create expedition: ID not set"
	ExpeditionPhoneNumberExists    = "Phone number already exists: %s"
	ExpeditionNotFound             = "expedition with id %s not found"

	// Success messages
	ExpeditionDeleteSuccess = "Successfully deleted Expedition"
)

// Helper functions for formatted error messages
func ExpeditionPhoneNumberExistsError(phoneNumber string) string {
	return fmt.Sprintf(ExpeditionPhoneNumberExists, phoneNumber)
}
