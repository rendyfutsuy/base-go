package constants

const (
	DefaultAuditor = "system"

	InvalidJsonInput   = "Invalid JSON input"
	InvalidJsonRequest = "Invalid JSON request"

	ContentType             = "application/json"
	FieldContentType        = "Content-Type"
	FieldContentDisposition = "Content-Disposition"

	ErrorJson = "Error decoding JSON : "

	ExcelContent = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	// TIME
	FormatDate                   = "2006-01-02"
	FormatDateTime               = "2006-01-02 15:04:05"
	FormatTimestamp              = "2006-01-02 15:04:05.999999999"
	FormatFullTimestamp          = "2006-01-02 15:04:05.000000 +0000 +0000"
	FormatFullTimestampGMT7      = "2006-01-02 15:04:05.999999999 -0700 -0700"
	FormatDateFileName           = "2006_01_02"
	FormatDateTimeFileName       = "2006_01_02_15_04_05"
	FormatDateTimeFileNameTicket = "060102_150405"
	FormatDateTimeString         = "2 January 2006 15:04:05"
	FormatDateString             = "2 January 2006"
	FormatDateStringMin          = "2 Jan 2006"

	//Status Mobile
	StatusSuccess = "success"
	StatusFailed  = "failed"
)
