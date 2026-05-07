package util

import (
	"errors"
	"strings"
	"time"
)

const (
	DefaultFormat        = "20060102"
	NowDefaultFormat     = "060102"
	FormatYYYYMM         = "200601"
	FormatYYYYMMDDHHMMSS = "20060102150405"
	StandardFormat       = "2006-01-02 15:04:05"
	StandardDayFormat    = "2006-01-02"
	PreciseFormat        = "2006-01-02 15:04:05.000"
	EmptyDateTime        = "0001-01-01T00:00:00Z"
	DefaultDateTime      = "0000-00-00 00:00:00"
	MaxDateTime          = "9999-12-31 23:59:59"
)

func TimeFormatDefault(t time.Time) string {
	return t.Format(DefaultFormat)
}

func TimeIntFormatYYYYMMDD(t int64) string {
	return time.Unix(t, 0).Format(DefaultFormat)
}

func TimeIntFormatYYMMDD(t int64) string {
	return time.Unix(t, 0).Format(NowDefaultFormat)
}

func TimeIntFormatYYMM(t int64) string {
	return time.Unix(t, 0).Format(FormatYYYYMM)
}

func NowFormatDYYMMDD() string {
	return time.Now().Format(NowDefaultFormat)
}

// 获取当前时间标准格式字符串。preciseMode：精确模式，为true精确到毫秒，false精确到秒
func NowStandardFormat(preciseMode bool) string {
	layout := StandardFormat
	if preciseMode {
		layout = PreciseFormat
	}
	return time.Now().Format(layout)
}

func NowStandardDayFormat() string {
	return time.Now().Format(StandardDayFormat)
}

func TimeStandardFormat(time time.Time, preciseMode bool) string {
	layout := StandardFormat
	if preciseMode {
		layout = PreciseFormat
	}
	return time.Format(layout)
}

// 转换 2006-01-02T15:04:05+08:00 格式时间为 2006-01-02 15:04:05
func FormatTime(t string) string {
	if t == "" {
		return ""
	}
	return strings.ReplaceAll(t[0:19], "T", " ")
}

func StrTimeToTimePtr(tm string) *time.Time {
	if len(tm) != 19 {
		return nil
	}
	t, err := time.ParseInLocation(StandardFormat, tm, time.Local)
	if err != nil {
		return nil
	}
	return &t
}

func StrTimeToTime(tm string) (time.Time, error) {
	if len(tm) != 19 {
		return time.Time{}, errors.New("时间格式不正确")
	}
	return time.ParseInLocation(StandardFormat, tm, time.Local)
}

func GetStartAndEndOfDay(t time.Time) (start, end time.Time) {
	// 获取当日开始时间
	start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	// 获取次日开始时间，并减去1纳秒得到当日结束时间
	end = start.Add(24 * time.Hour).Add(-1 * time.Nanosecond)
	return start, end
}
