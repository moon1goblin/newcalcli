package cmdshit_test

import (
	"calcli/cmdshit"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timeFromStrTestcases = []struct{
	Testname string

	Time_str string

	// making timestr would be a pain
	// so ill just verify with strings
	ExpectedRes string
	ExpectedWasOnlyDate bool
	ExpectedErr error
}{
	{"just a normal test", "12/12/12/12", "2025-12-12 12:12:00", false, nil},
	{"multiple delimiters", "11 --. 9 /// 00006", "2025-09-11 06:00:00", false, nil},
	{"empty string", "", "", false, cmdshit.ErrEmptyString},
	{"no month", "12", "", false, cmdshit.ErrNoDayAndMonth},
	{"invalid date", "blyat", "", false, cmdshit.ErrInvalidDateTime},
	{"zeros", "0 0", "", false, cmdshit.ErrInvalidDateTime},
	{"month literal first", "jan 25", "2025-01-25 00:00:00", true, nil},
	{"month literal second", "8 Mar 1:42", "2025-03-08 01:42:00", false, nil},
	{"whitespace first", " 9 10", "2025-10-09 00:00:00", true, nil},
	{"zeros in time", "9.08 00:00", "2025-08-09 00:00:00", false, nil},
}

func TestTimeFromStr(t *testing.T) {
	for _, tc := range timeFromStrTestcases {
		t.Run(tc.Testname, func(t *testing.T){
			assert := assert.New(t)
			res, onlydate, err := cmdshit.TimeFromStr(tc.Time_str)

			if res.Valid {
				assert.Equal(res.Time.Format(time.DateTime), tc.ExpectedRes)
				assert.Equal(onlydate, tc.ExpectedWasOnlyDate)
			}
			assert.ErrorIs(err, tc.ExpectedErr)
		})
	}
}
