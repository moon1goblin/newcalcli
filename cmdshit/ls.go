package cmdshit

import (
	"calcli/dbshit"
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
	events, err := dbshit.GetEventsInRange(
		begin_time,
		end_time,
		db_ptr,
	)
	if err != nil {
		return fmt.Errorf("lsAction error: %w", err)
	}

	// dereferencing nil ptr is ub i think haha
	if events == nil {
		return nil
	}

	for _, event := range *events {
		fmt.Println(
			event.Id,
			event.Name,
			event.Begin_time.String(),
			func() string {
				if event.End_time.Valid {
					return event.End_time.Time.String()
				}
				return "null"
			}(),
			event.Type,
		)
	}

	return nil
}
