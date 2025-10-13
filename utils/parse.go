package utils

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// Handle string formatted currency to float64
// Example:
// '1000' → 1000.000000
// '1,000' → 1000.000000
// '1.000' → 1000.000000
// '1000.50' → 1000.500000
// '1000,50' → 1000.500000
// '1,234.56' → 1234.560000
// '1.234,56' → 1234.560000
// '$2,500.42' → 2500.420000
// 'Rp 3.500,99' → 3500.990000
func ParseCurrency(input string) (float64, error) {
	// Check and store negative sign
	isNegative := false
	if strings.HasPrefix(input, "-") {
		isNegative = true
		input = input[1:]
	}

	// Strip all non-numeric, non-separator characters (keep digits, ',' and '.')
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == ',' || r == '.' {
			return r
		}
		return -1
	}, input)

	dotCount := strings.Count(cleaned, ".")
	commaCount := strings.Count(cleaned, ",")

	if dotCount == 0 && commaCount == 0 {
		val, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, err
		}
		if isNegative {
			val = -val
		}
		return val, nil
	}

	if dotCount > 0 && commaCount == 0 {
		lastDot := strings.LastIndex(cleaned, ".")
		if len(cleaned)-lastDot-1 == 2 {
			val, err := strconv.ParseFloat(cleaned, 64)
			if err != nil {
				return 0, err
			}
			if isNegative {
				val = -val
			}
			return val, nil
		}
		cleaned = strings.ReplaceAll(cleaned, ".", "")
		val, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, err
		}
		if isNegative {
			val = -val
		}
		return val, nil
	}

	if commaCount > 0 && dotCount == 0 {
		lastComma := strings.LastIndex(cleaned, ",")
		if len(cleaned)-lastComma-1 == 2 {
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
			val, err := strconv.ParseFloat(cleaned, 64)
			if err != nil {
				return 0, err
			}
			if isNegative {
				val = -val
			}
			return val, nil
		}
		cleaned = strings.ReplaceAll(cleaned, ",", "")
		val, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, err
		}
		if isNegative {
			val = -val
		}
		return val, nil
	}

	lastDot := strings.LastIndex(cleaned, ".")
	lastComma := strings.LastIndex(cleaned, ",")
	if lastComma > lastDot {
		cleaned = strings.ReplaceAll(cleaned, ".", "")
		cleaned = strings.ReplaceAll(cleaned, ",", ".")
	} else {
		cleaned = strings.ReplaceAll(cleaned, ",", "")
	}

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, err
	}
	if isNegative {
		val = -val
	}
	return val, nil
}

func ParseFloat(val interface{}) float64 {
	switch v := val.(type) {
	case string:
		res, _ := strconv.ParseFloat(strings.ReplaceAll(v, ",", ""), 64)
		return res
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1.0
		}
		return 0.0
	default:
		return 0.0
	}
}

func ParseBool(val interface{}) bool {
	switch v := val.(type) {
	case string:
		if strings.EqualFold(v, "yes") {
			return true
		}

		res, _ := strconv.ParseBool(v)
		return res
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	case bool:
		return v
	default:
		return false
	}
}

func ParseIsTBABool(val interface{}) bool {
	switch v := val.(type) {
	case string:
		if strings.ToLower(v) == "tba" {
			return true
		}

		res, _ := strconv.ParseBool(v)
		return res
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	case bool:
		return v
	default:
		return false
	}
}

func ParseStructToMap(s any) map[string]any {
	result := make(map[string]any)

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	// If pointer to struct, dereference
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Get json tag
		tag := field.Tag.Get("json")
		if tag == "-" || tag == "" {
			continue
		}

		// Handle tag with options (e.g., json:"name,omitempty")
		tagName := tag
		if commaIdx := indexComma(tag); commaIdx != -1 {
			tagName = tag[:commaIdx]
		}

		result[tagName] = value.Interface()
	}

	return result
}
