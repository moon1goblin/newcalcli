package cmdshit

import (
	"calcli/dbshit"
	"database/sql"
	"context"
	"errors"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_new *cli.Command = &cli.Command{
	Name: "new",
	// FIXME: flags think their values are only the first word :(
	// so like "-n im batman"'s value is im
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
	Action: newAction,
}

func newAction(ctx context.Context, cmd *cli.Command) error {
	// process dates somehow
	p_begin_time, err := processDate(cmd.String("begin"))
	if err != nil {
		return err
	}

	my_event := dbshit.Event{
		Name: cmd.String("name"),
		Begin_time: *p_begin_time,
	}

	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	found, err := my_event.Find(db_ptr)
	if err != nil {
		return err
	}
	if found {
		return errors.New("newAction error: event already exists")
	}

	if err := my_event.Push(db_ptr); err != nil {
		return err
	}
	return nil
}
