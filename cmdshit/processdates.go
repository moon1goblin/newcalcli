package cmdshit

import(
	"calcli/dbshit"
	"errors"
	"fmt"
)

var(
	ErrInvalidBeginEndCombo = errors.New("invalid begin/end combination")
	ErrEndBeforeBegin = errors.New("end before begin")
)

func ProcessDates(event_name_str, begin_datetime_str, end_datetime_str string) (*dbshit.Event, error) {
	p_end_time, onlydate_end, err_end := TimeFromStr(end_datetime_str)
	if errors.Is(err_end, ErrEmptyString) {
		// end can be null
	} else if err_end != nil {
		return nil, fmt.Errorf(
			"ProcessDates error with begin_datetime %s and end_datetime %s: %w: %w",
			begin_datetime_str,
			end_datetime_str,
			ErrInvalidDateTime,
			err_end,
		)
	}
	p_begin_time, onlydate_begin, err_begin := TimeFromStr(begin_datetime_str)
	if err_begin != nil {
		return nil, fmt.Errorf(
			"ProcessDates error with begin_datetime %s and end_datetime %s: %w: %w",
			begin_datetime_str,
			end_datetime_str,
			ErrInvalidDateTime,
			err_begin,
		)
	}

	// chronological order check
	if p_end_time.Valid && p_end_time.Time.Before(p_begin_time.Time) {
		return nil, fmt.Errorf(
			"ProcessDates error with begin_datetime %s and end_datetime %s: %w",
			begin_datetime_str,
			end_datetime_str,
			ErrEndBeforeBegin,
		)
	}

	var my_event_type dbshit.EventType

	// fuck this logic it was a pain to write
	if onlydate_begin && (!p_end_time.Valid || onlydate_end) {
		my_event_type = dbshit.FullDayEvent
	} else if !onlydate_begin && !p_end_time.Valid {
		my_event_type = dbshit.InstantEvent
	} else if !onlydate_begin && !onlydate_end {
		my_event_type = dbshit.WithDurationEvent
	} else {
		return nil, fmt.Errorf(
			"ProcessDates error with begin_datetime %s and end_datetime %s: %w",
			begin_datetime_str,
			end_datetime_str,
			ErrInvalidBeginEndCombo,
		)
	}

	return &dbshit.Event{
		Name: event_name_str,
		Begin_time: p_begin_time.Time,
		End_time: p_end_time,
		Type: my_event_type,
	}, nil
}
