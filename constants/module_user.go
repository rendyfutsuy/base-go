package constants

const (
	UserNotFound           = "User Not Found"
	UserEmailAlreadyExists = "User Email already exists"
	UserInvalid            = "User Not Found, please check to Customer Services..."
	UserIDNotFound         = "user user with id %s not found"
	
	// User import error messages
	UserEmailAlreadyExistsID    = "Email sudah terdaftar di database"
	UserUsernameAlreadyExistsID = "Username sudah terdaftar di database"
	UserNikAlreadyExistsID     = "NIK sudah terdaftar di database"
	
	// Excel import file errors
	UserImportExcelOpenFailed         = "failed to open Excel file"
	UserImportExcelNoSheets           = "Excel file has no sheets"
	UserImportExcelReadFailed         = "failed to read Excel file"
	UserImportExcelInsufficientRows   = "Excel file must have at least header row and one data row"
	UserImportFileNotFound            = "File tidak ditemukan. Gunakan field 'file' untuk upload Excel file"
	UserImportInvalidFileFormat       = "File harus berformat .xlsx atau .xls"
	UserImportFileOpenFailed          = "Gagal membuka file"
	UserImportTempFileCreateFailed    = "Gagal membuat file temporary"
	UserImportFileSaveFailed          = "Gagal menyimpan file"
	UserImportTemplateCreateFailed    = "Gagal membuat template"
	
	// Excel import validation errors
	UserImportRowInsufficientColumns  = "Row tidak memiliki cukup kolom (minimal 5 kolom: email, full_name, username, nik, role_name)"
	UserImportEmailRequired           = "email tidak boleh kosong"
	UserImportFullNameRequired        = "full_name tidak boleh kosong"
	UserImportUsernameRequired        = "username tidak boleh kosong"
	UserImportNikRequired             = "nik tidak boleh kosong"
	UserImportRoleNameRequired        = "role_name tidak boleh kosong"
	UserImportEmailInvalidFormat      = "Format email tidak valid"
	UserImportBatchDuplicationFailed  = "failed to check batch duplication"
	UserImportRoleNotFound            = "Role dengan nama '%s' tidak ditemukan"
	UserImportBatchCreateFailed       = "Error creating user in batch"
)
