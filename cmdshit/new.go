package cmdshit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

var Cmd_new *cli.Command = &cli.Command{
	Name: "new",
	// TODO: flags think their values are only the first word :(
	// so like "-n im batman"'s value is im
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "begin",
			Aliases: []string{"b"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "end",
			Aliases: []string{"e"},
			Required: false,
		},
		&cli.StringFlag{
			Name: "name",
			Aliases: []string{"n"},
			Required: true,
		},
	},
	Action: newAction,
}

var ErrEventAlreadyExists = errors.New("newAction error: event already exists")

func newAction(ctx context.Context, cmd *cli.Command) error {
	my_event, err := ProcessDates(cmd.String("name"), cmd.String("begin"), cmd.String("end"))
	if err != nil {
		return fmt.Errorf("newAction error: %w", err)
	}

	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)

	if found, err := my_event.Find(db_ptr); err != nil {
		return fmt.Errorf("newAction error: %w", err)
	} else if found {
		return ErrEventAlreadyExists
	}

	if err := my_event.Push(db_ptr); err != nil {
		return fmt.Errorf("newAction error: %w", err)
	}

	return nil
}
