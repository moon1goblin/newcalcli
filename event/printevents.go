package event

import (
	"strings"
)

func (event Event) String() string {
	var builder strings.Builder

	switch event.Type {
	case FullDayEvent:
	case InstantEvent:
		builder.WriteString(event.Begin_time.Format("15:04"))
		builder.WriteString(" ")
	case WithDurationEvent:
		builder.WriteString(event.Begin_time.Format("15:04"))
		builder.WriteString("-")
		builder.WriteString(event.End_time.Time.Format("15:04"))
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
	)

	// FIXME: fix multiple day events
	for _, event := range *events {
		if cur_date := event.Begin_time.Format("mon 2 jan"); cur_date != last_date {
			last_date = cur_date
			if first_run {
				first_run = false
			} else {
				builder.WriteString("\n")
			}
			builder.WriteString(cur_date)
			builder.WriteString("\n")
		}
		builder.WriteString(event.String())
		builder.WriteString("\n")
	}

	return builder.String()
}
