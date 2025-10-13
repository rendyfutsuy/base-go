package utils

import (
	"database/sql/driver"
	"encoding/json"
)

// NullInt64 is a custom nullable int64 type
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	if value == nil {
		ni.Int64, ni.Valid = 0, false
		return nil
	}
	ni.Valid = true
	if v, ok := value.(int64); ok {
		ni.Int64 = v
	} else if v, ok := value.(int32); ok {
		ni.Int64 = int64(v)
	} else {
		ni.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface for NullInt64
func (ni NullInt64) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int64, nil
}

// MarshalJSON handles JSON serialization for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON handles JSON deserialization for NullInt64
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	var v *int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		ni.Valid = false
		return nil
	}
	ni.Int64 = *v
	ni.Valid = true
	return nil
}
