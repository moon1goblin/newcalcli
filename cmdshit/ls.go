package cmdshit

import (
	"calcli/dbshit"
	"database/sql"
	"context"
	"time"
	"fmt"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

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
	var(
		begin_time *time.Time
		end_time *time.Time
		err error
	)

	if begin_time, err = processDate(cmd.String("begin")); err != nil && 
		err.Error() != "processDate error: empty string" {
		return err
	}
	if end_time, err = processDate(cmd.String("end")); err != nil && 
		err.Error() != "processDate error: empty string" {
		return err
	}

	// get sorted events in range [begin, end)
	events, err := dbshit.GetEventsInRange(
		begin_time,
		end_time,
		db_ptr,
	)
	if err != nil {
		return err
	}

	// dereferencing nil ptr is ub i think haha
	if events == nil {
		return nil
	}
	for _, event := range *events {
		fmt.Println(event.Id, event.Name, event.Begin_time)
	}

	return nil
}
