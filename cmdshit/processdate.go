package cmdshit

import (
	"strings"
	"strconv"
	"errors"
	"slices"
	"time"
)

var delimiters = []rune{':', ' ', '.', '/', '-'}

// TODO: make better error wrappers bruh

// can return nil if err
func processDate(time_str string) (*time.Time, error) {
	// empty string check
	if time_str == "" {
		return nil, errors.New("processDate error: empty string")
	}

	// day month hour minute
	datetimevalues := [4]int{}

	last_value_was_delimiter := false
	var builder strings.Builder
	var err error
	cur_datetime_value := 0

	proccessLastSlice := func() error {
		if !last_value_was_delimiter {
			if datetimevalues[cur_datetime_value], err = strconv.Atoi(builder.String()); err != nil {
				// TODO:
				// add ability to type months like jan or January instead of 1 here
				// so like if err check if its a valid month and have it that way
				return errors.New("processDate error: invalid datetime")
			}
			cur_datetime_value++
			builder.Reset()
		}
		return nil
	}

	for _, value := range time_str {
		if !slices.Contains(delimiters, value) {
			builder.WriteRune(value)
			last_value_was_delimiter = false
		} else {
			if err := proccessLastSlice(); err != nil {
				return nil, err
			}
			last_value_was_delimiter = true
		}
	}
	// process the last one
	if err := proccessLastSlice(); err != nil {
		return nil, err
	}

	// if no month and day we dont like that
	if datetimevalues[0] == 0 && datetimevalues[1] == 0 {
		return nil, errors.New("processDate error: date must have day and month")
	}

	datetime := time.Date(
		time.Now().Year(),
		// can i do this?
		time.Month(datetimevalues[1]),
		datetimevalues[0],
		datetimevalues[2],
		datetimevalues[3],
		0,
		0,
		time.Now().Location(),
	)

	return &datetime, nil
}
