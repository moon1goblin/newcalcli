package dbshit

import(
	"time"
	"errors"
)

// fuck sqlite datetimes

type TimeStr struct {
	datetimeval time.Time
	stringval string
}

func (t TimeStr) String() string {
	return t.stringval
}

func (t TimeStr) Time() time.Time {
	return t.datetimeval
}

func TimeStrFromStr(str string) (*TimeStr, error) {
	newtimeval, err := time.Parse(time.DateTime, str)
	if err != nil {
		return nil, errors.New("TimeStr error: failed to convert string to time.Time")
	}
	return &TimeStr{
		datetimeval: newtimeval,
		stringval: str,
	}, nil
}

func TimeStrFromTime(datetime time.Time) *TimeStr {
	return &TimeStr{
		datetimeval: datetime,
		stringval: datetime.Format(time.DateTime),
	}
}
