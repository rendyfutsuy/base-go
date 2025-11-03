package utils

import "math"

func ValidatePercentile(value float64) bool {
	temp := true

	if math.IsNaN(value) {
		temp = false
	} else {
		if value > 100 || value <= 0 {
			temp = false
		}
	}
	return temp
}
