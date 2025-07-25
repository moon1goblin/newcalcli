package dbshit

import (
	"errors"
	"fmt"
	"time"
)

// fuck sqlite datetimes

type TimeStr struct {
	datetimeval time.Time
	stringval string
	is_null bool
}

func (t TimeStr) String() *string {
	if t.is_null {
		return nil
	}
	return &t.stringval
}

func (t TimeStr) Time() *time.Time {
	if t.is_null {
		return nil
	}
	return &t.datetimeval
}

var ErrNullString = errors.New("nil pointer to string")

func TimeStrFromStr(str *string) (*TimeStr, error) {
	if str == nil {
		return nil, fmt.Errorf("TimeStr error: failed to convert string to time.Time: %w", ErrNullString)
	}
	newtimeval, err := time.Parse(time.DateTime, *str)
	if err != nil {
		return nil, fmt.Errorf("TimeStr error: failed to convert string to time.Time: %w", err)
	}
	return &TimeStr{
		datetimeval: newtimeval,
		stringval: *str,
	}, nil
}

func TimeStrFromTime(datetime time.Time) *TimeStr {
	return &TimeStr{
		datetimeval: datetime,
		stringval: datetime.Format(time.DateTime),
	}
}
