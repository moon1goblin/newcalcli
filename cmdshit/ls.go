package cmdshit

import (
	"calcli/event"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

// TODO: dislpay events prettier

var Cmd_ls *cli.Command = &cli.Command{
	Name: "ls",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "begin",
			Aliases: []string{"b"},
			Required: false,
		},
		&cli.StringFlag{
			Name: "end",
			Aliases: []string{"e"},
			Required: false,
		},
	},
	Action: lsAction,
}

func lsAction(ctx context.Context, cmd *cli.Command) error {
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	// process dates
	begin_time, _, err := TimeFromStr(cmd.String("begin"))
	if err != nil && !errors.Is(err, ErrEmptyString) {
		return fmt.Errorf("lsAction error on begin flag: %w", err)
	}
	end_time, _, err := TimeFromStr(cmd.String("begin"))
	if err != nil && !errors.Is(err, ErrEmptyString) {
		return fmt.Errorf("lsAction error on end flag: %w", err)
	}

	// get sorted events in range [begin, end)
	events, err := event.GetEventsInRange(
		begin_time,
		end_time,
		db_ptr,
	)
	if err != nil {
		return fmt.Errorf("lsAction error: %w", err)
	}

	fmt.Print(event.PrintEvents(events))

	return nil
}
