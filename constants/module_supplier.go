package constants

import "fmt"

const (
	SupplierContactTypeTelp  = "Telp"
	SupplierContactTypePhone = "Phone"
)

const (
	// Supplier validation errors
	SupplierInvalidRelationDateFormat = "invalid relation_date format, expected YYYY-MM-DD"
	SupplierCreateFailedIDNotSet      = "failed to create supplier: ID not set"
	SupplierIdentityNumberExists       = "Identity number already exists"
	SupplierPhoneNumberExists          = "Phone number already exists: %s"
	SupplierSubdistrictNotFound        = "Subdistrict not found"
	SupplierDistrictNotFound           = "District not found"
	SupplierCityNotFound               = "City not found"
	SupplierExpeditionNotFound         = "Expedition not found"
	SupplierDeliveryOptionNotFound     = "Delivery option not found"
	SupplierExpeditionPaidByNotFound   = "Expedition paid by not found"
	SupplierExpeditionCalculationNotFound = "Expedition calculation not found"

	// File upload errors
	SupplierFileOpenFailed        = "failed to open uploaded file: %s"
	SupplierFileReadFailed        = "failed to read uploaded file: %s"
	SupplierFileUploadFailed      = "failed to upload identity document: %s"

	// Success messages
	SupplierDeleteSuccess = "Successfully delete Supplier"
)

// Helper functions for formatted error messages
func SupplierPhoneNumberExistsError(phoneNumber string) string {
	return fmt.Sprintf(SupplierPhoneNumberExists, phoneNumber)
}

func SupplierFileOpenFailedError(err error) string {
	return fmt.Sprintf(SupplierFileOpenFailed, err.Error())
}

func SupplierFileReadFailedError(err error) string {
	return fmt.Sprintf(SupplierFileReadFailed, err.Error())
}

func SupplierFileUploadFailedError(err error) string {
	return fmt.Sprintf(SupplierFileUploadFailed, err.Error())
}
