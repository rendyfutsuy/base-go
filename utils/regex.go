package utils

import "regexp"

var EndorsedFacultativeRegex = regexp.MustCompile(`-E\d+$`)

func IsStringNumeric(s string) bool {
	match, _ := regexp.MatchString(`^\d+$`, s)
	return match
}
