package utils

import (
	"database/sql/driver"
	"encoding/json"
	"math"
)

// NullInt32 is a custom nullable int32 type
type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

// Scan implements the Scanner interface for NullInt32
func (ni *NullInt32) Scan(value interface{}) error {
	if value == nil {
		ni.Int32, ni.Valid = 0, false
		return nil
	}
	ni.Valid = true
	if v, ok := value.(int32); ok {
		ni.Int32 = v
	} else if v, ok := value.(int64); ok && v >= math.MinInt32 && v <= math.MaxInt32 {
		ni.Int32 = int32(v)
	} else {
		ni.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface for NullInt32
func (ni NullInt32) Value() (driver.Value, error) {
	if !ni.Valid {
		return nil, nil
	}
	return ni.Int32, nil
}

// MarshalJSON handles JSON serialization for NullInt32
func (ni NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}

// UnmarshalJSON handles JSON deserialization for NullInt32
func (ni *NullInt32) UnmarshalJSON(data []byte) error {
	var v *int32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		ni.Valid = false
		return nil
	}
	ni.Int32 = *v
	ni.Valid = true
	return nil
}
