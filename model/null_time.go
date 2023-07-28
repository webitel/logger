package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type NullTime time.Time

// - database/sql.Valuer
func (t NullTime) Value() (driver.Value, error) {
	if !t.IsZero() {
		return t.Time(), nil
	}
	return nil, nil
}

// IsZero value (?)
func (t *NullTime) IsZero() bool {
	if t != nil {
		return false
	}
	return true
}

// IsZero value (?)
func (t *NullTime) String() string {
	if t != nil {
		return t.Time().String()
	}
	return ""
}

// Time native value (!)
func (t *NullTime) Time() (stamp time.Time) {
	// v == time.Time{} // Zero(!)
	if t != nil {
		stamp = (time.Time)(*t) // shallowcopy(!)
	}
	return stamp
}

// Scan native value decoder function
func (t *NullTime) Scan(v interface{}) error {
	// scan: nullable (!)
	if v == nil {
		(t) = nil // Zero(!)
		return nil
	}
	switch v := v.(type) {
	//case int64:
	//	// +OK: timestamp
	//	if v == 0 {
	//		*(t) = Timestamp{} // Zero(!)
	//	} else {
	//		const timestampToUnix = (int64)(TimestampToUnix)
	//		switch TimestampToUnix {
	//		case time.Second:
	//			*(t) = (Timestamp)(time.Unix(v, 0))
	//		case time.Millisecond,
	//			time.Microsecond,
	//			time.Nanosecond:
	//			*(t) = (Timestamp)(time.Unix(v/timestampToUnix,
	//				v%timestampToUnix*(int64)(UnixToTimestamp),
	//			))
	//		default:
	//			panic(fmt.Errorf("timestamp: invalid precision %e", float64(UnixToTimestamp)))
	//		}
	//	}
	//	return nil
	//case []byte:
	//	return t.UnmarshalText(v)
	//case string:
	//	return t.UnmarshalText([]byte(v))
	case time.Time:
		// +OK: datetime
		(t) = (*NullTime)(&v) // shallowcopy
		return nil
	//case *time.Time:
	//	// +OK: datetime
	//	*(t) = (Timestamp)(*v) // shallowcopy
	//	return nil
	//case *Timestamp:
	//	if v != nil && t != v {
	//		*(t) = *(v) // shallowcopy
	//	}
	//	return nil
	default:
	}
	return fmt.Errorf("nulltime: convert %[1]T value %[1]v", v)
}
