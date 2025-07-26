package cmdshit_test

import (
	"calcli/cmdshit"
	"calcli/event"
	"testing"

	"github.com/stretchr/testify/assert"
)

var processDatesTestcases = []struct{
	Testname string

	Event_name_str string
	Begin_datetime_str string
	End_datetime_str string

	// ok for this test itd be 0 if event == nil
	ExpectedEventType event.EventType
	ExpectedErr error
}{
	{"full day 1 day", "", "12 12", "", event.FullDayEvent, nil},
	{"full day multiple days", "", "8 11", "10 11", event.FullDayEvent, nil},
	{"both empty", "", "", "", event.NilEvent, cmdshit.ErrEmptyString},
	{"instant", "", "9.08 15:00", "", event.InstantEvent, nil},
	{"duration", "", "05/08  008:5", "5-8 8:6", event.WithDurationEvent, nil},
	{"end before begin", "", "1 1 5:31", "1 1 5:30", event.NilEvent, cmdshit.ErrEndBeforeBegin},
	{"invalid combo", "", "3.04 15:15", "8.09", event.NilEvent, cmdshit.ErrInvalidBeginEndCombo},
}

func TestProccesDate(t *testing.T) {
	for _, tc := range processDatesTestcases {
		t.Run(tc.Testname, func(t *testing.T){
			assert := assert.New(t)
			event, err := cmdshit.ProcessDates(
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
