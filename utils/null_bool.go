package utils

import (
    "database/sql/driver"
    "encoding/json"
)

// NullBool is a custom nullable boolean type
type NullBool struct {
    Bool  bool
    Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface for NullBool
func (nb *NullBool) Scan(value interface{}) error {
    if value == nil {
        nb.Bool, nb.Valid = false, false
        return nil
    }
    nb.Valid = true
    switch v := value.(type) {
    case bool:
        nb.Bool = v
    case []byte:
        nb.Bool = string(v) == "true"
    case string:
        nb.Bool = v == "true"
    default:
        nb.Valid = false
        return nil
    }
    return nil
}

// Value implements the driver Valuer interface for NullBool
func (nb NullBool) Value() (driver.Value, error) {
    if !nb.Valid {
        return nil, nil
    }
    return nb.Bool, nil
}

// MarshalJSON handles JSON serialization for NullBool
func (nb NullBool) MarshalJSON() ([]byte, error) {
    if !nb.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(nb.Bool)
}

// UnmarshalJSON handles JSON deserialization for NullBool
func (nb *NullBool) UnmarshalJSON(data []byte) error {
    var b *bool
    if err := json.Unmarshal(data, &b); err != nil {
        return err
    }
    if b != nil {
        nb.Bool = *b
        nb.Valid = true
    } else {
        nb.Valid = false
    }
    return nil
}