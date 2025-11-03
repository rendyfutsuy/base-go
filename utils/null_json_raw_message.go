package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// NullJSONRawMessage is a custom nullable time type
type NullJSONRawMessage struct {
	JsonRawMessage json.RawMessage
	Valid          bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface for NullJSONRawMessage
func (nt *NullJSONRawMessage) Scan(value interface{}) error {
	if value == nil {
		nt.JsonRawMessage, nt.Valid = json.RawMessage{}, false
		return nil
	}
	nt.Valid = true
	t, ok := value.([]byte)
	if !ok {
		return errors.New("could not convert value to json.RawMessage")
	}
	result := json.RawMessage{}
	err := json.Unmarshal(t, &result)
	if err != nil {
		return err
	}
	nt.JsonRawMessage = result
	return nil
}

// Value implements the driver Valuer interface for NullJSONRawMessage
func (nt NullJSONRawMessage) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	value, err := nt.JsonRawMessage.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return value, nil
}

// MarshalJSON handles JSON serialization for NullJSONRawMessage
func (nt NullJSONRawMessage) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return nt.JsonRawMessage.MarshalJSON()
}

// UnmarshalJSON handles JSON deserialization for NullJSONRawMessage
func (nt *NullJSONRawMessage) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	if err := nt.JsonRawMessage.UnmarshalJSON(data); err != nil {
		return err
	}
	nt.Valid = true
	return nil
}
