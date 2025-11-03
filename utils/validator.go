package utils

import(
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func initValidator() *validator.Validate {
	return validator.New()
}