package cal_test

import (
	"github.com/moon1goblin/newcalcli/cal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var processDatesTestcases = []struct {
	Testname string

	Event_name_str     string
	Begin_datetime_str string
	End_datetime_str   string

	// ok for this test itd be 0 if event == nil
	ExpectedEventType cal.EventType
	ExpectedErr       error
}{
	{"full day 1 day", "", "12 12", "", cal.FullDayEvent, nil},
	{"full day multiple days", "", "8 11", "10 11", cal.FullDayEvent, nil},
	{"both empty", "", "", "", cal.NilEvent, cal.ErrEmptyString},
	{"instant", "", "9.08 15:00", "", cal.InstantEvent, nil},
	{"duration", "", "05/08  008:5", "5-8 8:6", cal.WithDurationEvent, nil},
	{"end before begin", "", "1 1 5:31", "1 1 5:30", cal.NilEvent, cal.ErrEndBeforeBegin},
	{"invalid combo", "", "3.04 15:15", "8.09", cal.NilEvent, cal.ErrInvalidBeginEndCombo},
}

func TestProccesDate(t *testing.T) {
	for _, tc := range processDatesTestcases {
		t.Run(tc.Testname, func(t *testing.T) {
			assert := assert.New(t)
			event, err := cal.ProcessDates(
				tc.Event_name_str,
				tc.Begin_datetime_str,
				tc.End_datetime_str,
			)

			assert.ErrorIs(err, tc.ExpectedErr)

			if event == nil {
				assert.Equal(int(tc.ExpectedEventType), 0)
			} else {
				assert.Equal(event.Type, tc.ExpectedEventType)
			}
		})
	}
}

var timeFromStrTestcases = []struct {
	Testname string

	Time_str string

	// making timestr would be a pain
	// so ill just verify with strings
	ExpectedRes         string
	ExpectedWasOnlyDate bool
	ExpectedErr         error
}{
	{"just a normal test", "12/12/12/12", "2025-12-12 12:12:00", false, nil},
	// FIXME: this test below
	{"ignore seconds and ms", "12/12/12/12/12/12", "2025-12-12 12:12:00", false, nil},
	{"multiple delimiters", "11 --. 9 /// 00006", "2025-09-11 06:00:00", false, nil},
	{"empty string", "", "", false, cal.ErrEmptyString},
	{"no month", "12", "", false, cal.ErrNoDayAndMonth},
	{"invalid date", "blyat", "", false, cal.ErrInvalidDateTime},
	{"zeros", "0 0", "", false, cal.ErrInvalidDateTime},
	{"month literal first", "jan 25", "2025-01-25 00:00:00", true, nil},
	{"month literal second", "8 Mar 1:42", "2025-03-08 01:42:00", false, nil},
	{"whitespace first", " 9 10", "2025-10-09 00:00:00", true, nil},
	{"zeros in time", "9.08 00:00", "2025-08-09 00:00:00", false, nil},
}

// FIXME: figure out how to run just one test from here
// because the function were testing here
// depends on the function we tested earlier in this file

func TestTimeFromStr(t *testing.T) {
	for _, tc := range timeFromStrTestcases {
		t.Run(tc.Testname, func(t *testing.T) {
			assert := assert.New(t)
			res, onlydate, err := cal.TimeFromStr(tc.Time_str)

			if res.Valid {
				assert.Equal(res.Time.Format(time.DateTime), tc.ExpectedRes)
				assert.Equal(onlydate, tc.ExpectedWasOnlyDate)
			}
			assert.ErrorIs(err, tc.ExpectedErr)
		})
	}
}
