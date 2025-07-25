package cmdshit_test

import (
	"calcli/cmdshit"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testcases = []struct{
	Testname string
	Time_str string
	// making timestr would be a pain
	// so ill just verify with strings
	ExpectedRes string
	ExpectedErr error
}{
	{"just a normal test", "12/12/12/12", "2025-12-12 12:12:00", nil},
	{"multiple delimiters", "11 --. 9 /// 00006", "2025-09-11 06:00:00", nil},
	{"empty string", "", "", cmdshit.ErrEmptyString},
	{"no month", "12", "", cmdshit.ErrNoDayAndMonth},
	{"invalid date", "blyat", "", cmdshit.ErrInvalidDateTime},
	{"zeros", "0 0", "", cmdshit.ErrInvalidDateTime},
	{"month literal first", "jan 25", "2025-01-25 00:00:00", nil},
	{"month literal second", "8 Mar 1:42", "2025-03-08 01:42:00", nil},
	{"whitespace first", " 9 10", "2025-10-09 00:00:00", nil},
}

func TestProccesDate(t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.Testname, func(t *testing.T){
			assert := assert.New(t)
			res, err := cmdshit.ProcessDate(tc.Time_str)

			if res != nil {
				assert.Equal(res.String(), tc.ExpectedRes)
			}
			assert.ErrorIs(err, tc.ExpectedErr)
		})
	}
}
