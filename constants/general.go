package constants

import "errors"

var (
	UrlAndErrorDigitalSignEmpyty = "url and error digital sign is empty"
	ErrorDataAlreadyExists       = errors.New("data already exists")
	TimezoneAsiaJakarta          = "Asia/Jakarta"
)

var (
	ErrForeignKeyViolation = errors.New("record is still referenced in another table")
)
