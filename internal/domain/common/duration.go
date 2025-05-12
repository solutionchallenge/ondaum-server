package common

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type DurationFormat string

const (
	DurationFormatSeconds = DurationFormat("3600sec")
	DurationFormatTime    = DurationFormat("02:04:05")
	DurationFormatMinutes = DurationFormat("60min")
)

type Duration time.Duration

func (d *Duration) Scan(value interface{}) error {
	if value == nil {
		*d = Duration(0)
		return nil
	}

	switch v := value.(type) {
	case int64:
		*d = Duration(v)
	case []byte:
		duration, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*d = Duration(duration)
	default:
		return utils.NewError("unsupported Scan type for Duration: %T", value)
	}
	return nil
}

func (d Duration) Value() (driver.Value, error) {
	return int64(d), nil
}

func (d Duration) ToDuration() time.Duration {
	return time.Duration(d)
}

func (d Duration) ToString(format DurationFormat) string {
	seconds := int64(d.ToDuration().Seconds())

	switch format {
	case DurationFormatSeconds:
		return fmt.Sprintf("%dsec", seconds)
	case DurationFormatTime:
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	case DurationFormatMinutes:
		minutes := (seconds + 59) / 60
		return fmt.Sprintf("%dmin", minutes)
	default:
		return fmt.Sprintf("%dsec", seconds)
	}
}

func NewDuration(d time.Duration) Duration {
	return Duration(d)
}

type NullableDuration struct {
	Duration Duration
	Valid    bool
}

func (nd *NullableDuration) Scan(value interface{}) error {
	if value == nil {
		nd.Duration, nd.Valid = Duration(0), false
		return nil
	}

	nd.Valid = true
	switch v := value.(type) {
	case int64:
		nd.Duration = Duration(v)
	case []byte:
		duration, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		nd.Duration = Duration(duration)
	default:
		return utils.NewError("unsupported Scan type for NullableDuration: %T", value)
	}
	return nil
}

func (nd NullableDuration) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}
	return int64(nd.Duration), nil
}

func (nd NullableDuration) ToDuration() time.Duration {
	return nd.Duration.ToDuration()
}

func (nd NullableDuration) ToString(format DurationFormat) string {
	if !nd.Valid {
		return ""
	}
	return nd.Duration.ToString(format)
}

func NewNullableDuration(d time.Duration) NullableDuration {
	if d == 0 {
		return NullableDuration{
			Duration: Duration(0),
			Valid:    false,
		}
	}
	return NullableDuration{
		Duration: NewDuration(d),
		Valid:    true,
	}
}
