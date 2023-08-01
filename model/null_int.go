package model

import (
	"database/sql/driver"
	"fmt"
)

type NullInt int

// - database/sql.Valuer
func (t NullInt) Value() (driver.Value, error) {
	if !t.IsZero() {
		return t.Int64(), nil
	}
	return nil, nil
}

// IsZero value (?)
func (t *NullInt) IsZero() bool {
	if t != nil {
		return false
	}
	return true
}

// IsZero value (?)
func (t *NullInt) Int() int {
	if t != nil {
		return (int)(*t)
	}
	return 0
}

// IsZero value (?)
func (t *NullInt) Int64() int64 {
	if t != nil {
		return (int64)(*t)
	}
	return 0
}

func NewNullInt(i int) *NullInt {
	return (*NullInt)(&i)
}

func (t *NullInt) Scan(v interface{}) error {
	// scan: nullable (!)
	if v == nil {
		t = nil // Zero(!)
		return nil
	}

	switch v := v.(type) {
	case int:
		// +OK: int
		*t = (NullInt)(v)
		return nil
	case int64:
		// +OK: int
		newVal := (int)(v)
		*t = (NullInt)(newVal) // shallowcopy
		return nil
	case int32:
		// +OK: int
		newVal := (int)(v)
		*t = (NullInt)(newVal) // shallowcopy
		return nil

	default:
	}
	return fmt.Errorf("int: convert %[1]T value %[1]v", v)
}
