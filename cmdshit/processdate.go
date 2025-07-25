package cmdshit

import (
	"calcli/dbshit"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

var delimiters = []rune{':', ' ', '.', '/', '-'}

var(
	ErrEmptyString = errors.New("processDate error: empty string")
	ErrInvalidDateTime = errors.New("processDate error: invalid datetime")
	ErrNoDayAndMonth = errors.New("processDate error: date must have day and month")
)

// can return nil if err
func processDate(time_str string) (*dbshit.TimeStr, error) {
	// empty string check
	if time_str == "" {
		return nil, ErrEmptyString
	}

	// day month hour minute
	datetimevalues := [4]int{}

	last_value_was_delimiter := false
	var builder strings.Builder
	var err error
	cur_datetime_value := 0

	proccessLastSlice := func() error {
		if last_value_was_delimiter {
			return nil
		}
		if datetimevalues[cur_datetime_value], err = strconv.Atoi(builder.String()); err != nil {
			// TODO:
			// add ability to type months like jan or January instead of 1 here
			// so like if err check if its a valid month and have it that way
			return fmt.Errorf("%w on string %s", ErrInvalidDateTime, time_str)
		}
		cur_datetime_value++
		builder.Reset()
		return nil
	}

	for _, value := range time_str {
		if !slices.Contains(delimiters, value) {
			builder.WriteRune(value)
			last_value_was_delimiter = false
			continue
		}
		if err := proccessLastSlice(); err != nil {
			return nil, err
		}
		last_value_was_delimiter = true
	}
	// process the last one
	if err := proccessLastSlice(); err != nil {
		return nil, err
	}

	// if no month and day we dont like that
	if cur_datetime_value <= 1  {
		return nil, fmt.Errorf("%w on string %s", ErrNoDayAndMonth, time_str)
	}

	datetime := time.Date(
		// TODO: add year somehow
		time.Now().Year(),
		time.Month(datetimevalues[1]),
		datetimevalues[0],
		datetimevalues[2],
		datetimevalues[3],
		// fuck seconds and ms
		0,
		0,
		time.Now().Location(),
	)

	return dbshit.TimeStrFromTime(datetime), nil
}
