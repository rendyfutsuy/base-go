package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/rendyfutsuy/base-go/constants"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return ValidateRequest(i, cv.Validator)
}

// getJSONFieldName returns the JSON field name for a given struct field
// register your custom message here
func ValidateRequest(req interface{}, valStruck *validator.Validate) error {
	// Register custom validator
	// registerCustomValidator(valStruck)

	if err := valStruck.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errorMessages []string

			for _, fe := range ve {
				fieldName := getJSONFieldName(req, fe.StructField())

				if fe.Tag() == "required" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
				}

				if fe.Tag() == "datetime" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s wrong from expected format", fieldName))
				}

				if fe.Tag() == "max" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("The characters you entered is too long. Please shorten it and try again"))
				}

				if fe.Tag() == "min" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("The %s Value You Enter is Insufficient", fe.Field()))
				}

				if fe.Tag() == "eqfield" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("The Confirmation value not same with the intended value"))
				}

				if fe.Tag() == "oneof" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s: Please Select One of Intended Value..", fieldName))
				}

				if fe.Tag() == "email" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("Your Email is not Valid.."))
				}

				if fe.Tag() == "emaildomain" {
					// get env variable
					emailScope := ConfigVars.String("email.validation-scope")

					// explore emailScope to array string
					emailScopes := strings.Split(emailScope, "|")

					// implode emailScopes
					emailScopesStr := strings.Join(emailScopes, ", ")

					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("Your Email is not Valid.., you must use `%s` domain", emailScopesStr))
				}

				if fe.Tag() == "nullableDate" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s: Date Not Valid (ex:2022-11-31)", fieldName))
				}

				if fe.Tag() == "required_if" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s: now required", fieldName))
				}

				if fe.Tag() == "password_uppercase" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s must be alphanumeric, minimum 8 characters, and all letters must be uppercase", fieldName))
				}

				if fe.Tag() == "uppercase_letters" {
					// Construct a human-friendly error message
					errorMessages = append(errorMessages, fmt.Sprintf("%s must have all alphabetic characters in uppercase", fieldName))
				}
			}
			// Join all error messages and return as a single string
			errorMessage := strings.Join(errorMessages, ", ")
			return errors.New(errorMessage)
		}
		// Fallback for any other errors
		return err
	}
	return nil
}

// getJSONFieldName returns the JSON field name for a given struct field
// shall not be called on other class by itself
func getJSONFieldName(val interface{}, fieldName string) string {
	// Get the field from the struct
	structField, found := reflect.TypeOf(val).Elem().FieldByName(fieldName)
	if !found {
		return strings.ToLower(fieldName) // Return the struct field name in lowercase as a fallback
	}

	// Get the label tag from the field
	labelTag := structField.Tag.Get("label")
	if labelTag != "" {
		return labelTag
	}

	// Get the JSON tag from the field
	jsonTag := structField.Tag.Get("json")
	if jsonTag == "" {
		return strings.ToLower(fieldName) // Return the struct field name in lowercase as a fallback
	}

	// The JSON tag might contain options like `json:"field,omitempty"`. We need only the field name.
	jsonFieldName := strings.Split(jsonTag, ",")[0]
	return jsonFieldName
}

// Add your Custom Validation
// Custom validation function for date format (nullable)
// shall not be called on other class by itself
func validateNullableDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" { // Skip validation if the field is empty
		return true
	}
	_, err := time.Parse(constants.FormatDate, dateStr)
	return err == nil
}

// Custom validation function for email domain
// shall not be called on other class by itself
func validateEmailDomain(fl validator.FieldLevel) bool {
	emailScope := ConfigVars.String("email.validation-scope")
	email := fl.Field().String()
	regexPattern := `^.*@(` + emailScope + `)$`
	regex := regexp.MustCompile(regexPattern)
	return regex.MatchString(email)
}

// Custom validation function for password
// Password must be alphanumeric, minimum 8 characters, and all letters must be uppercase
// shall not be called on other class by itself
func validatePasswordUppercase(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	// Check if all characters are alphanumeric and uppercase
	// Regex: ^[A-Z0-9]+$ means only uppercase letters and numbers
	regex := regexp.MustCompile(`^[A-Z0-9]+$`)
	return regex.MatchString(password)
}

// Custom validation function for uppercase letters
// Ensures all alphabetic characters in the string are uppercase
// shall not be called on other class by itself
func validateUppercaseLetters(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	if str == "" {
		return true // Skip validation if empty (use "required" tag for that)
	}

	// Check if all alphabetic characters are uppercase
	// Regex: ^[^a-z]*$ means no lowercase letters allowed
	for _, char := range str {
		if char >= 'a' && char <= 'z' {
			return false
		}
	}
	return true
}

// Helper function to check if a slice contains a particular value
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func RegisterCustomValidator(v *validator.Validate) {
	v.RegisterValidation("emaildomain", validateEmailDomain)
	v.RegisterValidation("nullableDate", validateNullableDate)
	v.RegisterValidation("password_uppercase", validatePasswordUppercase)
	v.RegisterValidation("uppercase_letters", validateUppercaseLetters)
}

// Helper function to get custom field name from struct tag
func getFieldName(req interface{}, field string) (fieldName string) {
	r := reflect.TypeOf(req)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	f, ok := r.FieldByName(field)
	if !ok {
		return strings.ToLower(field)
	}

	fieldName = f.Tag.Get("label")
	if fieldName == "" {
		fieldName = strings.ToLower(f.Name)
	}

	return fieldName
}
