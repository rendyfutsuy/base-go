package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NullStringArray is a custom nullable type for []string
type NullStringArray struct {
	Strings []string
	Valid   bool // Valid is true if Strings is not NULL
}

// Scan implements the Scanner interface for NullStringArray
func (nsa *NullStringArray) Scan(value interface{}) error {
	if value == nil {
		nsa.Strings, nsa.Valid = nil, false
		return nil
	}
	nsa.Valid = true
	switch v := value.(type) {
	case []byte:
		var temp interface{}
		if err := json.Unmarshal(v, &temp); err != nil {
			return fmt.Errorf("json unmarshal error: %w", err)
		}
		switch temp.(type) {
		case []interface{}:
			if err := json.Unmarshal(v, &nsa.Strings); err != nil {
				return fmt.Errorf("json unmarshal error: %w", err)
			}
		default:
			return fmt.Errorf("json unmarshal error: expected array of strings but got %T", temp)
		}
	default:
		return fmt.Errorf("unsupported Scan type for NullStringArray: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface for NullStringArray
func (nsa NullStringArray) Value() (driver.Value, error) {
	if !nsa.Valid {
		return nil, nil
	}
	return json.Marshal(nsa.Strings)
}

// MarshalJSON handles JSON serialization for NullStringArray
func (nsa NullStringArray) MarshalJSON() ([]byte, error) {
	if !nsa.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nsa.Strings)
}

// UnmarshalJSON handles JSON deserialization for NullStringArray
func (nsa *NullStringArray) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nsa.Strings, nsa.Valid = nil, false
		return nil
	}
	if err := json.Unmarshal(data, &nsa.Strings); err != nil {
		return err
	}
	nsa.Valid = true
	return nil
}

// safe return array value, even if its out of bound

func SafeGetElement[T comparable](slice []T, index int) T {
	var defValue T
	if index < 0 || index >= len(slice) {
		return defValue
	}

	return slice[index]
}
