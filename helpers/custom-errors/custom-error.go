package custom_errors

import "strings"

type (
	CustomError struct {
		ErrorMesages []string `json:"error_mesages"`
		ArrayErrors  []error  `json:"array_errors"`
	}
)

func (err *CustomError) Error() string {
	var errorText string
	var tempErrorArray []string = err.ErrorMesages
	for _, v := range err.ArrayErrors {
		tempErrorArray = append(tempErrorArray, v.Error())
	}

	errorText = strings.Join(tempErrorArray, ";")

	return errorText
}

func (err *CustomError) AppendError() {
	for _, v := range err.ArrayErrors {
		err.ErrorMesages = append(err.ErrorMesages, v.Error())
	}

	return
}
