package utils

import (
	"database/sql/driver"
	"errors"
	"time"
)

// NullTime is a custom nullable time type
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface for NullTime
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Valid = true
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("could not convert value to time.Time")
	}
	nt.Time = t
	return nil
}

// Value implements the driver Valuer interface for NullTime
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON handles JSON serialization for NullTime
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return nt.Time.MarshalJSON()
}

// UnmarshalJSON handles JSON deserialization for NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	if err := nt.Time.UnmarshalJSON(data); err != nil {
		return err
	}
	nt.Valid = true
	return nil
}
