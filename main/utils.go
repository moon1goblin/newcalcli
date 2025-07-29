package main

import (
	"errors"
	"fmt"
	"github.com/moon1goblin/newcalcli/cal"
	"golang.org/x/term"
	"os"
	"strings"
)

// TODO: center events to the right edge of their time

// accepts SORTED array
func PrintEvents(events *[]cal.Event) string {
	if events == nil {
		return ""
	}
	var (
		builder   strings.Builder
		last_date string
		first_run bool = true
		// multiple_day_full_day_events = list.New()
	)

	for _, cur_event := range *events {
		// if cur_event.Type == cal.FullDayEvent {
		// 	multiple_day_full_day_events.PushBack(cur_event)
		// }
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
		builder.WriteString(cur_event.String())
		builder.WriteString("\n")
	}

	return builder.String()
}

// ------------------------------------------------------------------------------

var ErrYNPrompt = errors.New("y/n prompt error")

// so basically anything other than y or Y is false, good design hello
func ConfirmYNPrompt() (bool, error) {
	// stty into raw mode so we dont have to press enter
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return false, fmt.Errorf("ConfirmYNPrompt error setting terminal into raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	input := make([]byte, 1)
	if bytes_read, err := os.Stdin.Read(input); err != nil || bytes_read != 1 {
		return false, fmt.Errorf("ConfirmYNPrompt error reading char from stdin: %w", err)
	}

	fmt.Printf("%c", input[0])

	if input[0] == 'y' || input[0] == 'Y' {
		return true, nil
	} else {
		return false, nil
	}
}

func NumberSelectPrompt() (byte, error) {
	input := make([]byte, 2)
	if bytes_read, err := os.Stdin.Read(input); err != nil || bytes_read > 2 {
		return 0, fmt.Errorf("ConfirmYNPrompt error reading char from stdin: %w", err)
	}
	input[0] -= 48 // This is also retarded btw, but who cares
	input[1] -= 48
	// fmt.Printf("\n%d %d", input[0], input[1])
	// fmt.Printf("\n%d", input[0]*10+input[1])

	return input[0]*10 + input[1], nil // I know this is retarded but whatever
}
