package cmdshit

import (
	"context"
	"calcli/dbshit"
	"database/sql"

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
	my_event := dbshit.Event{
		Name: cmd.String("name"),
		Begin_time: cmd.String("begin"),
	}

	// take the db_ptr out of the context (again idk wtf that is)
	db_ptr := ctx.Value("db_ptr").(*sql.DB)
	return my_event.Push(db_ptr)
}
