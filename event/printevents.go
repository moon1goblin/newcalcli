package event

import (
	"strings"
	"time"
	"container/list"
)

func isSameDate(lhs, rhs time.Time) bool {
	ly, lm, ld := lhs.Date()
	ry, rm, rd := rhs.Date()
	return ld == rd && lm == rm && ly == ry
}

func (event Event) String(withdate bool) string {
	var builder strings.Builder

	var begin_format string
	var end_format string

	if withdate {
		begin_format = "01.02.2006 15:04"
	} else {
		begin_format = "15:04"
	}
	if event.End_time.Valid && isSameDate(event.Begin_time, event.End_time.Time) {
		end_format = "15:04"
	} else {
		end_format = begin_format
	}

	switch event.Type {
	case FullDayEvent:
	case InstantEvent:
		builder.WriteString(event.Begin_time.Format(begin_format))
		builder.WriteString(" ")
	case WithDurationEvent:
		builder.WriteString(event.Begin_time.Format(begin_format))
		builder.WriteString("-")
		builder.WriteString(event.End_time.Time.Format(end_format))
		builder.WriteString(" ")
	}

	builder.WriteString(event.Name)

	return builder.String()
}

// TODO: center events to the right edge of their time

// accepts SORTED array
func PrintEvents(events *[]Event) string {
	if events == nil {
		return ""
	}
	var(
		builder strings.Builder
		last_date string
		first_run bool = true
		multiple_day_full_day_events = list.New()
	)

	for _, cur_event := range *events {
		if cur_event.Type == FullDayEvent {
			multiple_day_full_day_events.PushBack(cur_event)
		}
		if cur_date := cur_event.Begin_time.Format("Mon 2 Jan"); cur_date != last_date {
			last_date = cur_date
			if first_run {
				first_run = false
			} else {
				builder.WriteString("\n")
			}
			builder.WriteString(cur_date)
			builder.WriteString("\n")

			// FIXME: multiple day event printing
			// // so we can have multiple day full day events duh
			// for elem := multiple_day_full_day_events.Front(); elem != nil; elem = elem.Next() {
			// 	if cur_fd_event, ok := elem.Value.(*Event); ok {
			// 		// !(cur < fulld) is same as fulldayevent >= cur_event i hope
			// 		if !cur_event.Begin_time.Before(cur_fd_event.Begin_time) {
			// 			builder.WriteString(cur_fd_event.String(false))
			// 		}
			// 		if cur_event.Begin_time.Before(cur_fd_event.Begin_time) {
			// 	} // it should be ok?
			// }
		}
		builder.WriteString(cur_event.String(false))
		builder.WriteString("\n")
	}

	return builder.String()
}
