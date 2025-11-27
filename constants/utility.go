package constants

const (
	DefaultAuditor         = "system"
	ErrorUUIDNotRecognized = "Requested param is not recognized" // Error when a requested parameter is not recognized.
	ErrorUUIDIsEmpty       = "Requested shipyard is not exists"  // Error when a requested shipyard does not exist.

	InvalidJsonInput   = "Invalid JSON input"
	InvalidJsonRequest = "Invalid JSON request"

	ContentType             = "application/json"
	FieldContentType        = "Content-Type"
	FieldContentDisposition = "Content-Disposition"

	ErrorJson = "Error decoding JSON : "

	ExcelContent = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

// ExcelContentDisposition returns Content-Disposition header value for Excel file download
func ExcelContentDisposition(filename string) string {
	return "attachment; filename=" + filename
}

const (
	// TIME
	FormatDate                           = "2006-01-02"
	FormatDateTime                       = "2006-01-02 15:04:05"
	FormatDateTime12H                    = "02-01-2006 03:04:05 PM"
	FormatDateTime12HFull                = "02 January 2006 03:04:05 PM"
	FormatTimezone                       = "2006-01-02T15:04:05.999Z07:00"
	FormatDateTimeISO8601                = "2006-01-02T15:04:05Z07:00" // ISO 8601 format with timezone offset
	FormatTimestamp                      = "2006-01-02 15:04:05.999999999"
	FormatFullTimestamp                  = "2006-01-02 15:04:05.000000 +0000 +0000"
	FormatFullTimestampGMT7              = "2006-01-02 15:04:05.999999999 -0700 -0700"
	FormatFullTimeStampType2             = "2006-01-02 00:00:00 +0000 +0000"
	FormatFullTimeStampType3             = "2006-01-02T15:04:05Z"
	FormatDateFileName                   = "2006_01_02"
	FormatDateTimeFileName               = "2006_01_02_15_04_05"
	FormatDateTimeFileNameTicket         = "060102_150405"
	FormatDateTimeString                 = "2 January 2006 15:04:05"
	FormatDateString                     = "2 January 2006"
	FormatDateStringZeroPaddedDay        = "02 January 2006"
	FormatDateStringZeroPaddedDayMin     = "02-Jan-06"
	FormatDateStringZeroPaddedDayMinYear = "02-Jan-2006"
	FormatDateStringZeroPaddedDayMonth   = "02-01-2006"
	FormatDateStringMin                  = "2 Jan 2006"
	FormatDateStringMinPadded            = "02 Jan 2006"
	FormatDateMMMDDYYYYZeroPaddedDay     = "Jan 02, 2006"
	FormatTimeConfig                     = "format.time"
	FormatFullDayDateTime                = "Monday, 02 January 2006 15:04:05"
	TimeEmpty                            = "0001-01-01T00:00:00Z"
	DateEmpty                            = "0001-01-01"
	FormatMonthYear                      = "January 2006"

	// Excel TIME
	FormatExcelDateTime    = "2006-01-02T15:04:05.999Z07:00"
	FormatExcelDateReverse = "02/01/2006"
	ExcelInvalidDate       = "01/01/0001"

	//Status Mobile
	StatusSuccess = "success"
	StatusFailed  = "failed"

	// sql message error
	SQLErrorQueryRow      = "QueryRow scan error: %v"
	SQLErrorQueryDatabase = "Error querying database:"
	SQLErrorScanRow       = "Error scanning row:"

	// file
	FileSizeTooLarge = "One or more files are too large"
	MultiFormEmpty   = "Failed to parse multipart form"

	// Number
	EPSILON = 1e-9

	ErrMalformedUserContext = "Could not process user data from context"
)
