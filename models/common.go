package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeWrapper struct {
	time.Time
}

func (t TimeWrapper) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		z := fmt.Sprintf("\"%s\"", "0000-00-00 00:00:00")
		return []byte(z), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.DateTime))
	return []byte(formatted), nil
}

func (t TimeWrapper) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *TimeWrapper) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = TimeWrapper{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *TimeWrapper) GetTime() time.Time {
	return t.Time
}
