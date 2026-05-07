package until

import "time"

func StrTimeToTimePtr(tm string) *time.Time {
	if len(tm) != 19 {
		return nil
	}
	t, err := time.ParseInLocation(time.DateTime, tm, time.Local)
	if err != nil {
		return nil
	}
	return &t
}
