package utils

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString is a custom nullable string type
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface for NullString
func (ns *NullBool) Scan(value interface{}) error {
	if value == nil {
		ns.Bool, ns.Valid = false, false
		return nil
	}
	ns.Valid = true
	if b, ok := value.([]byte); ok {
		ns.Bool = string(b) == "true"
	}
	return nil
}

// Value implements the driver Valuer interface for NullBool
func (ns NullBool) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.Bool, nil
}

// MarshalJSON handles JSON serialization for NullBool
func (ns NullBool) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Bool)
}

// UnmarshalJSON handles JSON deserialization for NullBool
func (ns *NullBool) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	ns.Bool = b
	ns.Valid = true
	return nil
}
