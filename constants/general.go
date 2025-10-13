package constants

import "errors"

// const (
// 	ITA23Percentage           = 2.0
// 	VATPercentage             = 2.2
// 	BrokerageInclusivePercent = 1.022
// )

var (
	UrlAndErrorDigitalSignEmpyty = "url and error digital sign is empty"
	ErrorDataAlreadyExists       = errors.New("data already exists")
	TimezoneAsiaJakarta          = "Asia/Jakarta"
)

var (
	ErrForeignKeyViolation = errors.New("record is still referenced in another table")
)
