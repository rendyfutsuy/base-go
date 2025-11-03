package utils

import (
	"database/sql/driver"
	"encoding/json"
)

// NullFloat64 is a custom nullable float64 type
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface for NullFloat64
func (ni *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		ni.Float64, ni.Valid = 0, false
		return nil
	}
	ni.Valid = true
	if v, ok := value.(float64); ok {
		ni.Float64 = v
	} else {
		ni.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface for NullFloat64
func (ni NullFloat64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Float64, nil
}

// MarshalJSON handles JSON serialization for NullFloat64
func (ni NullFloat64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Float64)
}

// UnmarshalJSON handles JSON deserialization for NullFloat64
func (ni *NullFloat64) UnmarshalJSON(data []byte) error {
	var v *float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		ni.Valid = false
		return nil
	}
	ni.Float64 = *v
	ni.Valid = true
	return nil
}
