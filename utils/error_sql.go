package utils

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
)

func HumanizePQError(err error) error {
	if err == nil {
		return nil
	}

	if pqErr, ok := err.(*pq.Error); ok {

		switch pqErr.Code {
		case "23505":

			consDetail := ExtractDetailFromConstraint(pqErr.Constraint)
			if consDetail != nil {
				return fmt.Errorf("%s %s already exists", consDetail["model"], consDetail["field"])
			}

			return constants.ErrorDataAlreadyExists
		case "23503":
			return fmt.Errorf("Oops! It seems that there is a foreign key constraint violation.")
		case "23502":
			return fmt.Errorf("Oops! It seems that a required field is missing.")
		case "23514":
			return fmt.Errorf("Oops! It seems that a check constraint is violated.")
		case "22001":
			return fmt.Errorf("The characters you entered is too long. Please shorten it and try again")
		default:
			return fmt.Errorf("An unexpected database error occurred")
		}
	}

	return err
}

// Helper function to extract the field name from the detail message
func ExtractFieldFromDetail(detail string) string {
	// Example detail message: "Key (name)=(The 1) already exists."
	start := strings.Index(detail, "(")
	end := strings.Index(detail, ")")
	if start != -1 && end != -1 && start < end {
		return detail[start+1 : end]
	}
	return "data"
}

// Helper function to extract detail from the constraint message
// in map type, model, and field from format "unique_class_name"
func ExtractDetailFromConstraint(constraint string) (res map[string]string) {
	// Validate constraint format
	if strings.Count(constraint, "_") != 2 {
		return nil
	}

	// Split the constraint into parts
	parts := strings.Split(constraint, "_")
	if len(parts) != 3 {
		return nil
	}

	// Assign parts to respective variables
	typeString := parts[0]
	modelString := parts[1]
	fieldString := parts[2]

	// Create and return the result map
	res = map[string]string{
		"type":  typeString,
		"model": modelString,
		"field": fieldString,
	}

	return res
}
