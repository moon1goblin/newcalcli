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

var(
	ErrEmptyString = errors.New("empty string")
	ErrNoDayAndMonth = errors.New("date must have day and month")
)

// TODO: add event types (fullday, instant, withduration)

// can return nil if err
// bool true when there was only a date and no time
func TimeFromStr(time_str string) (*dbshit.TimeStr, bool, error) {
	// empty string check
	if time_str == "" {
		return nil, false, fmt.Errorf("TimeFromStr error: %w", ErrEmptyString)
	}

	// day month hour minute
	datetimevalues := [4]int{}

	last_value_was_delimiter := true
	var builder strings.Builder
	cur_datetime_value := 0
	monthwasfirst := false

	proccessLastSlice := func() error {
		if last_value_was_delimiter {
			return nil
		}
		var (
			val int
			is_in_map bool
			err error
		)
		if val, is_in_map = months[builder.String()]; is_in_map {
			if cur_datetime_value == 0 {
				monthwasfirst = true
			}
		} else if val, err = strconv.Atoi(builder.String()); 
			err != nil || 
			cur_datetime_value <= 1 && val == 0 {
			return fmt.Errorf("TimeFromStr error on string %s: %w", time_str, ErrInvalidDateTime)
		}

		datetimevalues[cur_datetime_value] = val
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
			return nil, false, err
		}
		last_value_was_delimiter = true
	}
	// process the last one
	if err := proccessLastSlice(); err != nil {
		return nil, false, err
	}

	only_date_no_time := false

	// if no month and day we dont like that
	if cur_datetime_value <= 1  {
		return nil, false, fmt.Errorf("TimeFromStr error on string %s: %w", time_str, ErrNoDayAndMonth)
	} else if cur_datetime_value == 2 {
		only_date_no_time = true
	}

	if monthwasfirst {
		datetimevalues[0], datetimevalues[1] = datetimevalues[1], datetimevalues[0]
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

	return dbshit.TimeStrFromTime(datetime), only_date_no_time, nil
}

var delimiters = []rune{':', ' ', '.', '/', '-'}

// did this manually so i have more control := (= are tears)
// edit: it was so ugly i had to put it at the bottom of the file
var months = map[string]int {
	"January": 1,
	"Jan": 1,
	"january": 1,
	"jan": 1,

	"February": 2,
	"Feb": 2,
	"february": 2,
	"feb": 2,

	"March": 3,
	"Mar": 3,
	"march": 3,
	"mar": 3,

	"April": 4,
	"Apr": 4,
	"april": 4,
	"apr": 4,

	"May": 5,
	"may": 5,

	"June": 6,
	"Jun": 6,
	"june": 6,
	"jun": 6,

	"July": 7,
	"Jul": 7,
	"july": 7,
	"jul": 7,

	"August": 8,
	"Aug": 8,
	"august": 8,
	"aug": 8,

	"September": 9,
	"Sep": 9,
	"september": 9,
	"sep": 9,

	"October": 10,
	"Oct": 10,
	"october": 10,
	"oct": 10,

	"November": 11,
	"Nov": 11,
	"november": 11,
	"nov": 11,

	"December": 12,
	"Dec": 12,
	"december": 12,
	"dec": 12,
}
