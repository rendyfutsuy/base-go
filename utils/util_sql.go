package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
)

func SanitizeTableName(tableName string) string {
	tableName = strings.Replace(tableName, string([]byte{0}), "", -1)
	return `"` + strings.Replace(tableName, `"`, `""`, -1) + `"`
}

// StringToUUID converts a string to a UUID.
func StringToUUID(request interface{}) (result uuid.UUID, err error) {
	// Try to assert the request as a string.
	stringUUID, ok := request.(string)
	if !ok {
		// If it's not a string, try to assert it as a UUID.
		uuid, ok := request.(uuid.UUID)
		if !ok {
			// If it's neither, return an error.
			err = errors.New(constants.ErrorUUIDNotRecognized)
			return
		}
		// If it's a UUID, return it.
		result = uuid
	} else {
		// If it's a string, try to parse it as a UUID.
		uuid, parseErr := uuid.Parse(stringUUID)
		if parseErr != nil {
			// If the parsing fails, return an error.
			err = fmt.Errorf("requested param is string")
			return
		}
		// If the parsing succeeds, return the parsed UUID.
		result = uuid
	}
	return
}

func SearchingMonthToNumber(searchText, fieldMonthQuery string) string {
	// Map for month names to integers
	var monthMap = map[string]int{
		"january": 1, "february": 2, "march": 3, "april": 4, "may": 5, "june": 6,
		"july": 7, "august": 8, "september": 9, "october": 10, "november": 11, "december": 12,
	}

	var monthMapIndonesian = map[string]int{
		"januari": 1, "februari": 2, "maret": 3, "april": 4, "mei": 5, "juni": 6,
		"juli": 7, "agustus": 8, "september": 9, "oktober": 10, "november": 11, "desember": 12,
	}

	// Check if searchText contains a month name and replace it with the corresponding integer value
	for month, value := range monthMap {
		if strings.Contains(month, strings.ToLower(searchText)) {
			searchText = fmt.Sprintf("%d", value)
			break
		}
	}

	for month, value := range monthMapIndonesian {
		if strings.Contains(searchText, month) {
			searchText = fmt.Sprintf("%d", value)
			break
		}
	}

	// Check if searchText is a number
	if _, err := strconv.Atoi(searchText); err == nil {
		return " AND " + fieldMonthQuery + " = '" + searchText + "'"
	}

	return ""
}

func CleanString(input string) string {
	// Remove extra spaces and trim leading/trailing spaces
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(input, " "))
}
