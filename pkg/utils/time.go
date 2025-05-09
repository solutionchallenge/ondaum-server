package utils

import (
	"fmt"
	"time"
)

const (
	TIME_FORMAT_ISO8601 = "2006-01-02T15:04:05Z"
	TIME_FORMAT_DATE    = "2006-01-02"
)

func FormatTimezoneOffset(t time.Time) string {
	_, offset := t.Zone()
	hours := offset / 3600
	return fmt.Sprintf("UTC%+03d:00", hours)
}
