package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NullString is a custom nullable string type
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	switch v := value.(type) {
	case []byte:
		ns.String = string(v)
	case string:
		ns.String = v
	default:
		return fmt.Errorf("unsupported Scan type for NullString: %T", value)
	}
	return nil
}

// Value implements the driver Valuer interface for NullString
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON handles JSON serialization for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON handles JSON deserialization for NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ns.String = str
	ns.Valid = true
	return nil
}
