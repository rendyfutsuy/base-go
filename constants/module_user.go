package constants

const (
	UserNotFound                       = "User Not Found"
	UserEmailAlreadyExists             = "User Email already exists"
	UserInvalid                        = "User Not Found, please check to Customer Services..."
	UserIDNotFound                     = "user with id %s not found"
	UserPermissionModuleName           = "Users"
	UserPermissionNameCreate           = "Create"
	UserPermissionNameDelete           = "Delete"
	UserRestrictedPermissionGroupError = "Permission groups 'Create' and 'Delete' in module 'User' can only be assigned to the Super Admin role"

	// User import error messages
	UserEmailAlreadyExistsID    = "Email is already registered in the database"
	UserUsernameAlreadyExistsID = "Username is already registered in the database"
	UserNikAlreadyExistsID      = "NIK is already registered in the database"

	// Excel import file errors
	UserImportExcelOpenFailed       = "failed to open Excel file"
	UserImportExcelNoSheets         = "Excel file has no sheets"
	UserImportExcelReadFailed       = "failed to read Excel file"
	UserImportExcelInsufficientRows = "Excel file must have at least header row and one data row"
	UserImportFileNotFound          = "File not found. Use the 'file' field to upload the Excel file"
	UserImportInvalidFileFormat     = "File must be in .xlsx or .xls format"
	UserImportFileOpenFailed        = "Failed to open file"
	UserImportTempFileCreateFailed  = "Failed to create temporary file"
	UserImportFileSaveFailed        = "Failed to save file"
	UserImportFailedPartial         = "Failed to import some rows"
	UserImportTemplateCreateFailed  = "Failed to create template"

	// Excel import validation errors
	UserImportRowInsufficientColumns = "Row does not have enough columns (minimum 5 columns: email, full_name, username, nik, role_name)"
	UserImportEmailRequired          = "email cannot be empty"
	UserImportFullNameRequired       = "full_name cannot be empty"
	UserImportUsernameRequired       = "username cannot be empty"
	UserImportNikRequired            = "nik cannot be empty"
	UserImportRoleNameRequired       = "role_name cannot be empty"
	UserImportEmailInvalidFormat     = "Invalid email format"
	UserImportBatchDuplicationFailed = "failed to check batch duplication"
	UserImportRoleNotFound           = "Role with name '%s' was not found"
	UserImportBatchCreateFailed      = "Error creating user in batch"
	UserRoleNotFound                 = "Role not found"
	UserCannotDelete                 = "User cannot be deleted because deletable is false"
)
