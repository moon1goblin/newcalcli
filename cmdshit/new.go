package cmdshit

import (
	"calcli/dbshit"
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_new *cli.Command = &cli.Command{
	Name: "new",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "begin",
			Aliases: []string{"b"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "name",
			Aliases: []string{"n"},
			Required: true,
		},
	},
	Action: newFunc,
}

func newFunc(ctx context.Context, cmd *cli.Command) error {
	// process dates somehow
	// create event here
	p_begin_time, err := processDate(cmd.String("begin"))
	if err != nil {
		return err
	}
	my_event := dbshit.Event{
		Name: cmd.String("name"),
		Begin_time: p_begin_time,
	}
	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)
	return my_event.Push(db_ptr)
}

var delimiters = []rune{':', ' ', '.', '/', '-'}

func processDate(time_str string) (time.Time, error) {
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
				return time.Time{}, err
			}
			last_value_was_delimiter = true
		}
	}
	// process the last one
	if err := proccessLastSlice(); err != nil {
		return time.Time{}, err
	}

	// if no month and day we dont like that
	if datetimevalues[0] == 0 && datetimevalues[1] == 0 {
		return time.Time{}, errors.New("processDate error: date must have day and month")
	}

	return time.Date(
		time.Now().Year(),
		// can i do this?
		time.Month(datetimevalues[1]),
		datetimevalues[0],
		datetimevalues[2],
		datetimevalues[3],
		0,
		0,
		time.Now().Location(),
	), nil
}
